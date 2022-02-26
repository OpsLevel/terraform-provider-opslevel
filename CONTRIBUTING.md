# Contributing

1. [About this document](#about-this-document)
3. [Getting the code](#getting-the-code)
4. [Setting up an environment](#setting-up-an-environment)
5. [Local development](#local-development)
7. [Submitting a Pull Request](#submitting-a-pull-request)

## About this document

This document is a guide intended for folks interested in contributing to `terraform-provider-opslevel`. Below, we document the process by which members of the community should create issues and submit pull requests (PRs) in this repository. It is not intended as a guide for using `terraform`, and it assumes a certain level of familiarity with Golang and Terraform concepts. This guide assumes you are using macOS and are comfortable with the command line.

If you're new to Golang development or contributing to open-source software, we encourage you to read this document from start to finish.

## Proposing a change

This project is what it is today because community members like you have opened issues, provided feedback, and contributed to the knowledge loop for the entire communtiy. Whether you are a seasoned open source contributor or a first-time committer, we welcome and encourage you to contribute code, documentation, ideas, or problem statements to this project.

### Defining the problem

If you have an idea for a new feature or if you've discovered a bug, the first step is to open an issue. Please check the list of [open issues](https://github.com/OpsLevel/terraform-provider-opslevel/issues) before creating a new one. If you find a relevant issue, please add a comment to the open issue instead of creating a new one.

> **Note:** All community-contributed Pull Requests _must_ be associated with an open issue. If you submit a Pull Request that does not pertain to an open issue, you will be asked to create an issue describing the problem before the Pull Request can be reviewed.

### Submitting a change

If an issue is appropriately well scoped and describes a beneficial change to the codebase, then anyone may submit a Pull Request to implement the functionality described in the issue. See the sections below on how to do this.

The maintainers will add a `good first issue` label if an issue is suitable for a first-time contributor. This label often means that the required code change is small or a net-new addition that does not impact existing functionality. You can see the list of currently open issues on the [Contribute](https://github.com/OpsLevel/terraform-provider-opslevel/contribute) page.

## Getting the code

### Installing git

You will need `git` in order to download and modify the source code. On macOS, the best way to download git is to just install [Xcode](https://developer.apple.com/support/xcode/).

### External contributors

If you are not a member of the `OpsLevel` GitHub organization, you can contribute by forking the repository. For a detailed overview on forking, check out the [GitHub docs on forking](https://help.github.com/en/articles/fork-a-repo). In short, you will need to:

1. fork the repository
2. clone your fork locally
3. check out a new branch for your proposed changes
4. push changes to your fork
5. open a pull request from your forked repository

### OpsLevel contributors

If you are a member of the `OpsLevel` GitHub organization, you will have push access to the repo. Rather than forking to make your changes, just clone the repository, check out a new branch, and push directly to that branch.

## Setting up an environment

There are some tools that will be helpful to you in developing locally. While this is the list relevant for development in this repository, many of these tools are used commonly across open-source python projects.

### Tools

- []

## Local Development

### Installation

First make sure you have working [golang development environment](https://learn.gopherguides.com/courses/preparing-your-environment-for-go-development) setup. Also make sure you have the latest version of `terraform` [installed.](https://learn.hashicorp.com/tutorials/terraform/install-cli)

### Setup a Terraform workspace

We have a `workspace` folder in the repository that can be used as a place to play around with terraform.

After any code change you can just run the following in the `workspace` directory to build and pull in the latest provider code

```sh
./build
./init
```

To `terraform plan`

```sh
./plan
```

To `terraform apply`

```sh
./apply
```

Feel free to investigate the scripts to see what they are doing if you need to use some other terraform commands or adjust

## Submitting a Pull Request

OpsLevel provides a CI environment to test changes through Github Actions. For example, if you submit a pull request to the repo, GitHub will trigger automated code checks and tests upon approval from an OpsLevel maintainer.

A maintainer will review your PR. They may suggest code revision for style or clarity, or request that you add unit or integration test(s). These are good things! We believe that, with a little bit of help, anyone can contribute high-quality code.
- First time contributors should note code checks + unit tests require a maintainer to approve.

Once all tests are passing and your PR has been approved, a maintainer will merge your changes into the active development branch. And that's it! Happy developing :tada:
