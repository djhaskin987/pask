Why Pask?
=========

It's PAckages and tASKs.

Pask is given a list of Pask packages and is told either to install the packages
or run a specific task which is associated with all the installed packages.
Packages are installed relative to a project (build or deploy) root.

Here's the real power of pask: not only will the packages be installed
in the order that they appear in the list, but will also run the tasks
associated with the packages in the order that they appear in the given list.


Imagine running ``pask run deploy``, ``pask run blue-green``, or ``pask run
compile``, and the task associated with each package runs in order.

Well, at least *I* think it's pretty cool. :)

