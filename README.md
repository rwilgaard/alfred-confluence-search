# Confluence Search

A workflow for finding pages in your Confluence installation.

## Installation
* [Download the latest release](https://github.com/rwilgaard/alfred-confluence-search/releases)
* Open the downloaded file in Finder.
* If running on macOS Catalina or later, you _**MUST**_ add Alfred to the list of security exceptions for running unsigned software. See [this guide](https://github.com/deanishe/awgo/wiki/Catalina) for instructions on how to do this.

## Features
* Search for pages.
* Filter search by Confluence Space.
* Workflow auto update

## Keywords

You can change the default 'cs' keyword in the User configuration.
* With `cs` you can search for pages in Confluence. The default ↩ action will open the highlighted page in your browser.

## Filters

The following filters can be used in your query to filter your search:
* Filter pages by Space using `@space` syntax. Alfred will let you search and select a Space when you type `@` and complete the query.


## Actions

The following actions can be used on a highlighted page:
* `↩` will open the page in your browser.
* `⌘` + `↩` will copy the URL for the page to the clipboard.
