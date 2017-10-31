Quickstart
==========

1. Create a pask package called ``a`` with files in it::

    mkdir pkga
    cd pkga
    mkdir -p pask/tasks
    echo 'echo compile task A' > pask/tasks/mytask
    mkdir files
    touch files/a
    tar cJf ../a-1.0.tar.xz $(ls -1A)
    cd ..

2. Do the same thing to create a package called ``b``::

    mkdir pkgb
    cd pkgb
    mkdir -p pask/tasks
    echo 'echo compile task B' > pask/tasks/mytask
    mkdir files
    touch files/b
    tar cJf ../b-1.0.tar.xz $(ls -1A)
    cd ..

3. Write out a ``pask/spec.yml`` file relative to the root of your project::

      packages:
        - name: a
          version: 1.0
          location: file://<path-to-where-a-was-made>/a-1.0.tar.xz
        - name: b
          version: 1.3
          location: file://<path-to-where-b-was-made>/b-1.3.tar.xz

   Pask also accepts HTTP and HTTPS URLs.

4. Have pask install the contents of the packages::

       pask install
       find .

5. Run ``pask run mytask``::

       pask run mytask

