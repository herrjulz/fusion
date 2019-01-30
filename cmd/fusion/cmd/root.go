package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/JulzDiverse/fusion"
	"github.com/JulzDiverse/fusion/bpm"
	"github.com/JulzDiverse/fusion/docker/entrypoint"
	"github.com/JulzDiverse/fusion/kube"
	"github.com/JulzDiverse/fusion/monit"
	"github.com/JulzDiverse/fusion/packaging"
	"github.com/JulzDiverse/fusion/spec"
	"github.com/JulzDiverse/fusion/templates"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fuse",
	Short: "translate kubernetes specs and dockerfiles into a bosh release",
	Long:  `it does stuff`,
	Run:   fuse,
}

func init() {
	rootCmd.Flags().StringP("config", "c", "", "Path to the fusion config file")
}

func fuse(cmd *cobra.Command, args []string) {
	/*
		0. Parse Config (hard coded for now)
		1. Create Spec
		2. Create Templates
		3. Create BPM
		4. Create Monit-File
		5. Packaging
	*/

	conf := fusion.Config{
		BoshReleaseDir:     "test-release",
		Dockerfile:         "integration/Dockerfile",
		KubeDeploymentSpec: "integration/kube-deployment.yml",
		BinaryDownloadURL:  "www",
	}

	kubeSpecRaw, err := ioutil.ReadFile(conf.KubeDeploymentSpec)
	if err != nil {
		panic(err)
	}
	objectGroup := kube.ParseKubeObjects(kubeSpecRaw)

	//1. Create Spec
	specConverter := spec.Converter{}
	spec, err := specConverter.Convert(objectGroup)
	if err != nil {
		panic(err)
	}

	jobPath := filepath.Join(conf.BoshReleaseDir, "jobs", spec.Name)
	templatePath := filepath.Join(jobPath, "templates")
	err = os.MkdirAll(templatePath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	specfile, err := yaml.Marshal(spec)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filepath.Join(jobPath, "spec"), specfile, os.ModePerm)
	if err != nil {
		panic(err)
	}

	//2. Create Templates
	templateConverter := templates.Converter{}
	templates, err := templateConverter.Convert(objectGroup)
	if err != nil {
		panic(err)
	}

	for filename, contents := range templates {
		err = ioutil.WriteFile(filepath.Join(templatePath, filename), contents, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	//3. Create BPM
	entrypointParser := entrypoint.Parser{}
	bpmer := bpm.BPMer{entrypointParser}

	dockerfile, err := ioutil.ReadFile(conf.Dockerfile)
	if err != nil {
		panic(err)
	}

	bpmYamlObj, err := bpmer.ToBpm(spec.Name, bytes.NewBuffer(dockerfile))
	if err != nil {
		panic(err)
	}

	bpmYaml, err := yaml.Marshal(bpmYamlObj)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filepath.Join(templatePath, "bpm.yml.erb"), bpmYaml, os.ModePerm)
	if err != nil {
		panic(err)
	}

	//4. Create Monitfile
	monitfile := monit.Create("opi")
	err = ioutil.WriteFile(filepath.Join(jobPath, "monit"), []byte(monitfile), os.ModePerm)
	if err != nil {
		panic(err)
	}

	//5. Create Packaging
	packagingPath := filepath.Join(conf.BoshReleaseDir, "packages", spec.Name)
	err = os.MkdirAll(packagingPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	packageSpecFile := packaging.CreateSpec(spec.Name, "opi")
	packagingFile := packaging.CreateScript("opi", conf.BinaryDownloadURL)

	err = ioutil.WriteFile(filepath.Join(packagingPath, "spec"), []byte(packageSpecFile), os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filepath.Join(packagingPath, "packaging"), []byte(packagingFile), os.ModePerm)
	if err != nil {
		panic(err)
	}

	fmt.Println("Kubernetes Objects and Dockerfile successfully converted to bosh-release files!")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
