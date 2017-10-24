# circleci-runner

Use circleci-runner to download and execute your application binaries
directly from CircleCI and **never worry about distributing new versions of your
apps and services to your systems again**. circleci-runner always runs the
latest successful build, using your project's stored build artifacts.

New build? Don't worry. It also polls for new artifacts and uses best practices
to [gracefully shutdown][disposability] your app and execute the newly built
binary. Give it a shot!

[disposability]: https://12factor.net/disposability

## Quickstart

1. Setup your project to [store build artifacts][docs] of your application binary.

2. [Download the latest release binary][binary], [build from source][source], or
   [use Docker][docker] to get circleci-runner on your machine.

3. Go to your account dashboard on CircleCI and create a new
   [Personal API Token](https://circleci.com/account/api).

4. Export some environment variables in your terminal to configure your project:

   ```sh
   export CIRCLECI_TOKEN="your API token"
   export CIRCLECI_ACCOUNT="your user or team"
   export CIRCLECI_REPOSITORY="your repo name"
   ```

5. Finally, download and execute your project!

   ```sh
   circleci-runner <flags you want to pass to your app>
   ```

[docs]: https://circleci.com/docs/2.0/artifacts/#uploading-artifacts "Storing and Accessing Build Artifacts on CircleCI"
[binary]: https://github.com/hashrabbit/circleci-runner/releases "circleci-runner releases"
[source]: #building-from-source "How to build circleci-runner from source"
[docker]: #running-with-docker "How to run circleci-runner with Docker"

## Usage

You can configure your CircleCI credentials and control various settings
through environment variables. circleci-runner expects an executable artifact
(e.g., a native binary, shell script, or other executable). By default, it will
attempt to download and execute the first discovered artifact.

Pass arguments to your application binary on the command-line as regular
arguments to circleci-runner. They will be passed to your artifact like normal.

### Environment Variables

The following environment variables are available:

- `CIRCLECI_TOKEN` (**required**) - Your CircleCI API token. To generate a
  token, visit the *API Permissions* tab of your CircleCI project settings page
  and click *Create Token*. For additional security, limit its scope to
  *Build Artifacts*.

- `CIRCLECI_ACCOUNT` (**required**) - The team or personal account username that
  your project belongs to.

- `CIRCLECI_REPOSITORY` (**required**) - The repository name of your application.

- `CIRCLECI_BRANCH` (*optional*) - Limit builds to a specific branch (e.g., `master`).

- `CIRCLECI_ARTIFACT` (*optional*) - The path to the artifact to run, see the
  CircleCI docs regarding [build artifacts][build-artifacts] (defaults to the
  first found build artifact).

- `CIRCLECI_DEBUG` (*optional*) - Enable debug logging.

[build-artifacts]: https://circleci.com/docs/2.0/artifacts/ "Storing and Accessing Build Artifacts"

## Building From Source

You can build circleci-runner from source by having a [Go toolchain][toolchain]
installed and executing the following command:

```sh
go get github.com/hashrabbit/circleci-runner
```

Afterward, you should have circleci-runner in your `$PATH` and the source code
in your `$GOPATH`.

[toolchain]: https://golang.org/doc/install "Getting Started With Go"

## License

circleci-runner is licensed under the [Apache License, Version 2.0](LICENSE.md).
Read a summary of this license's permissions, conditions, and limitations
[here](https://choosealicense.com/licenses/apache-2.0/).
