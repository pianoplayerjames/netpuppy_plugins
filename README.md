# Requirements
- https://github.com/pianoplayerjames/netsquirrel

# Install Plugins
- when you're running netpuppy you can install a plugin by typing ```install``` to run the installer and then the name of the plugin. you can either type ```name_of_plugin``` or ```name_of_plugin.go```. After installing you may need to restart netpuppy for it to work.

# Changing plugin store
- in netpuppy you can change the plugin store by editing plugins/install.go and changing the github repo to your own url. for example: ```https://raw.githubusercontent.com/<repo_username>/<repo_name>/<branch>```

make sure all the go plugins are in the root directory of your repo.
