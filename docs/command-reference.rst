Command Reference for Pask
==========================

Root CLI Reference
------------------

The basic idea is that a pask package is simply an archive that contains files
and tasks. These archives can be installed and their tasks can be run in order
of a list given to pask when it starts its run. All packages will be installed
and run relative to the project root directory, presumably of a build or
deployment.

If you were to run ``pask help`` today you might see this output::

    Install files and run tasks. Make your build
    less painful and more fun!

    Usage:
      pask [command]

    Available Commands:
      help        Help about any command
      install     Install or update a package or packages
      run         Run a packaged task

    Flags:
      -p, --base string     Base project path (default "./")
      -c, --config string   config file (default is $HOME/.pask)
      -h, --help            help for pask
      -s, --spec string     Pask spec file

    Use "pask [command] --help" for more information about a command.

Explanation of root command options:

* ``-p BASE``, ``--base BASE``: Specify base project path. This is the path
  relative to which pask will try to install files and/or run tasks. By
  default this is the present working directory (``./``).
* ``-c FILE``, ``--config FILE``: Specify config file for pask. This file
  is in `TOML`_ format. By default it can be found at ``~/.pask``. This
  file currently only allows the operator to specify values for the
  base project path and spec file location relative to that path, using the 
  ``base`` and ``spec`` file keys. Arguably configuring these is useless, but
  the option to use a configuration file is there to facilitate growth of
  pask and its configuration values for future versions of pask.
* ``-s SPEC``, ``--spec SPEC``: Specify pask spec file location. By default
  this location is ``<base>/pask/spec.yml``. This file is intended to be
  committed to project source code and specifies, under the ``package`` key,
  the list of packages that pask will use for its operations. Each entry in
  the list under the ``package`` key itself has at least three keys: ``name``,
  ``version``, and ``location``. Example::

      packages:
        - name: a
          version: 1.0
          location: file://<path-to-where-a-was-made>/a-1.0.tar.xz
        - name: b
          version: 1.3
          location: file://<path-to-where-b-was-made>/b-1.3.tar.xz

.. _TOML: https://github.com/toml-lang/toml

CLI Reference for ``pask install``
----------------------------------

The command ``pask install`` has no command-specific options in the current
version. 

When `pask install` is run, this list of packages is installed in order of the
list. Each package is simply a `XZ-compressed`_ `tarball`_ is downloaded and
unpacked relative to the build project root. If any files are found in the
archive under the folder `./pask`, those files are installed under
`./pask/packages/<package-name>/<package-version>/` relative to the build
project root.

.. _XZ-compressed: https://en.wikipedia.org/wiki/Xz
.. _tarball: https://en.wikipedia.org/wiki/Tar_(computing)

CLI Reference for ``pask run`` 
----------------------------------

The command ``pask run`` has no command-specific options in the current
version. 

When `pask run <task>` is run, the file
`./pask/packages/<package-name>/<package-version>/tasks/<task>` is assumed to
be executable. It is run with no arguments for each package in the list found
in the file `./pask/spec.yml` relative to the build project root. Any
environment variables set when pask is run are passed through. If, for any of
the packages in the list of the spec file, the named task does not exist, is
not executable, or exits abnormally, an error is printed and pask stops. The
task for each package in the list are run in the order that the packages found
in the list.

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

.. _Degasolv: https://degasolv.readthedocs.io
