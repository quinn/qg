function config(args) {
    return args
}

function addGenerator(fileData, config) {
    // check for ending newline
    if (!fileData.endsWith('\n')) {
        fileData += '\n'
    }

    return fileData + '  - name: ' + config.name + '\n'
}
