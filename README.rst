===============================================================================
cidre-bottle: building a packaged cidre webapp binary
===============================================================================

.. image:: https://godoc.org/github.com/yuin/cidre-bottle?status.svg
    :target: http://godoc.org/github.com/yuin/cidre-bottle

|

cidre-bottle provides an easy way to build a fat binary for web applications using the `cidre <https://github.com/yuin/cidre/>`_ webframework.

----------------------------------------------------------------
Installation
----------------------------------------------------------------

.. code-block:: bash
   
   go get github.com/jteeuwen/go-bindata/...
   go get github.com/elazarl/go-bindata-assetfs/...
   go get github.com/yuin/cidre-bottle

----------------------------------------------------------------
Usage
----------------------------------------------------------------

directory structure::
    
    +- assets
       +- templates
          +- page.tpl
          +- layout.tpl
       +- statics
          +- css
             +- app.css
          +- js
             +- app.js
          +- img
             +- logo.png

Runs a `go-bindata` command.

.. code-block:: bash
    
    go-bindata assets/...

Modify cidre app codes.

.. code-block:: go

    app := cidre.NewApp(appConfig)
    // Set a renderer with a go-bindata support
    app.Hooks.Add("setup", func(w http.ResponseWriter, r *http.Request, data interface{}) {
        app.Renderer := bottle.NewHtmlTemplateRenderer(app.Renderer, Asset, AssetDir)
    })

    root := app.MountPoint("/")
    // Serve static files
    bottle.Static(root, "statics", "statics", "assets/statics", Asset, AssetDir)

----------------------------------------------------------------
License
----------------------------------------------------------------
MIT

----------------------------------------------------------------
Author
----------------------------------------------------------------
Yusuke Inuzuka
