Detailed Explanation
====================

The basic idea is that a pask package is simply an archive that contains files
and tasks. These archives can be installed and their tasks can be run in order
of a list given to pask when it starts its run. All packages will be installed
and run relative to the project root directory, presumably of a build or
deployment.

There are two subcommands, `install` and `run`.

Pask takes a list of packages in file relative to the
root of the build in YAML (or JSON) format called `./pask/spec.yml`. Pask will
look at the `packages` key in that file for a list,
with each entry in that list containing `name`, `version`, and `location`
keys for the packages to be installed and/or for the packages for which to
run tasks. Example::

      packages:
        - name: a
          version: 1.0
          location: file://<path-to-where-a-was-made>/a-1.0.tar.xz
        - name: b
          version: 1.3
          location: file://<path-to-where-b-was-made>/b-1.3.tar.xz

Pask Install
------------

When `pask install` is run, this list of packages is installed in order of the
list. Each package is simply a `XZ-compressed`_ `tarball`_ is downloaded and unpacked relative to the build project
root. The build project root can be changed using Pask's  If any files are found in the archive under the folder `./pask`, those
files are installed under `./pask/packages/<package-name>/<package-version>/`
relative to the build project root.

.. _XZ-compressed: https://en.wikipedia.org/wiki/Xz
.. _tarball: https://en.wikipedia.org/wiki/Tar_(computing)

Pask Run
--------

When `pask run <task>` is run, the file
`./pask/packages/<package-name>/<package-version>/tasks/<task>` is run with no
arguments for each package in the list found in the file `./pask/spec.yml`
relative to the build project root. Any environment variables set when pask
is run are passed through. If, for any of the packages in the list of the spec
file, the named task does not exist, is not executable, or exits abnormally, an
error is printed and pask stops. The task for each package in the list are run
in the order that the packages found in the list.

Compatibility with Degasolv
---------------------------

Pask was built to work natively with `Degasolv`_,
version 1.10.0 or greater. By using Degasolv to resolve dependencies between
pask packages, you get a complete package management system for your builds.

When you run ``degasolv resolve-locations``, Degasolv will print out a list of
packages in order of "dependance". Packages with many dependants appear first
in the list, while packages which have many dependencies appear last in the
list. 

The operator can create a valid spec for Pask by using Degasolv's
``--output-format`` CLI option, like so::

    cd <build-project-dir>
    mkdir -p ./pask
    degasolv resolve-locations --output-format json > pask/spec.yml

The Pask spec can then serve as Degasolv's lock file for the build project.
Then, Pask would run tasks associated with these tasks in order of dependance
(more dependers = installed first, more dependencies =
installed last), and will also run the tasks associated with those packages
in order of dependance.

.. _Degasolv: degasolv.readthedocs.io
