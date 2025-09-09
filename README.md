# devnote
> This tool is currently in development and not yet ready for use. First release is currently targeted for October 2025, subject to change.

<br>

**devnote** is a simple free tool that makes it easy to keep developer diaries. Keep track of all the bugs you've squashed, features you've built, and lessons you've learned, all in beautiful markdown format!

- Designed to keep your notes connected with your codebase through git integration and customisable templates.
- Perfect for students documenting their code, developers trying out a new language, or anybody wanting to log their projects development history through structured notes.

<br>

## Installation
- This project is currently in active development, and does not yet have any releases avaliable for installation.
- Can be built from source, however core features may be broken or incomplete.


<br>

## Basic Usage

### Initialise

To initialise devnotes in your current directory, run:

```
devnote init
```

- This creates the `.devnote` directory inside your project root, which will hold your local config, templates, and save data
- To allow for different behavior per project, each devnote instance has its own config file and set of templates
- For ease of use, your global config file and any global templates will be copied over to your new devnote instance
- If you wish for these to not be copied, use the `--clean` flag. This will create the default config and an empty templates folder

```
devnote init --clean
```

- A new folder is also created to store your notes. By default, this folder will be called `devnotes` and be held in your project root
- If you wish for the notes folder to be stored elsewhere, a custom path can be specified using the `--path` flag:

```
devnote init --path /home/Desktop/my-project-notes
```

### New

To create a new devnote, use the `new` command. It will automatically open in your preferred text editor

```
devnote new
```

- By default, the new devnote will use the currently selected template. To use a specific template, use the `--template` flag and pass in the name the desired template:

```
devnote new --template bug-fix
```

- You can use the `--blank` flag if you do not wish to use a template and want an empty file:

```
devnote new --blank
```

- If no file name is specified, the filename will be the timestamp of when the note was created. A name can be specified using the `--name` flag:

```
devnote new --name fixed-regression
```

- The devnote can be saved outside of the notes folder by specifying the desired path of the new note using the `--path` flag:

```
devnote new --path /home/Desktop
```

- These flags can be combined. `--path` should define the directory the file will be saved to. `--name` defines the name
```
devnote new --path /home/Desktop --name desktop-note
```


> **Note:** Saving a note outside of the notes folder will mean it does not show up when running `search` or `list`

### Viewing and Editing Notes

- To list all of your notes:
```
devnote list
```

- To view a specific note
```
devnote view <note-name>
```

- To search through your notes
```
devnote search <keyword>
```

- To edit an existing note in your configured editor
```
devnote edit <note-name>
```

### Configuration

- To view your current configuration:
```
devnote config
```
> **Note:** this command also validates your config file. If any errors are found, they will be output to the console instead

- To view a specfic setting:

```
devnote config get editor
```
This will output your currently selected text editor to terminal

- To set a configuration value:

```
devnote config set editor vim
```

- To perform any of these actions on the global config, use the `--global` flag
```
devnote config --global
devnote config set --global editor vim
```

<br>

> **Note:** Both the global can local configs can be edited manually. However, this approach is reccomended to avoid formatting issues

<br>

## Templating

Templating is a core part of devnote. Each time you create a new note, it will be pre-filled with text defined in the currently selected template

Templates are stored in the `.devnote/templates` directory

If you wish to use dyanmic data in the note, such as showing your latest git commits or the current time, you can use the following placeholders in your template

- `{project}` - The name of your project
- `{date}` - The current date
- `{time}` - The current time
- `{timestamp}` - The current timestamp (date + time)
- `{commits}` - Shows all git commits since the last devnote
- `{branch}` - The current working branch
- `{author-name}` - Your current configured git name
- `{author-email}` - Your current configured git email

<br>

## Advanced Usage

Since devnote creates standard markdown files, it plays nicely with the broader ecosystem:

- **Publish your dev diary online** - Link your notes folder with static website generators, like Jekyll or Hugo, to build an online presense for your project
- **Sync across devices and contributors** -
- **Integrate with documentation**

<br>

## Contributions

devnote is closed to contributions while in development. After first release, contributions will be welcome and encouraged! I will be responding to issues and merge requests here.

<br>

## Support

I'm so happy if you have found devnote useful, a lot of time and love went into making it :)

I’m a student building free and open source tools for fun, if you find my work helpful, feel free to support me on Ko-Fi — or just star a repo, share a project, or drop some feedback. It all means a lot!

<br>

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/G2G81GQB6Y)
