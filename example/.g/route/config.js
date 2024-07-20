/**
 * @typedef {Object} ConfigObject
 * @property {string} method
 * @property {string} path
 * @property {string} routeFilename
 * @property {string} viewFilename
 * @property {string} funcName */

/**
 * @param {Object} options
 * @param {string} options.method
 * @param {string} options.path
 * @returns {ConfigObject} */
function config({ method, path }) {
    method = method.toUpperCase()

    // remove first char of path if it is '/'
    if (path.startsWith('/')) {
        rpath = path.slice(1)
    }

    const parts = rpath.split('/')
    const routeFilename = parts.map((part) => {
        if (part.startsWith(':')) {
            return part.replace(':', `$`)
        }

        return part
    }).join('.')

    const viewFilename = parts.filter((part) => !part.startsWith(':')).join('-')

    const funcName = parts.filter((part) => !part.startsWith(':')).map((part) => {
        return part.charAt(0).toUpperCase() + part.slice(1)
    }).join('')

    return { method, path, routeFilename, viewFilename, funcName }
}


/**
 * @param {string} fileData
 * @param {ConfigObject} config
 * @returns {string} */
function addRoute(fileData, config) {
    const out = ''
    for (const line of fileData.split('\n')) {
        if (line.endsWith('/* insert new routes here */')) {
            out += `e.${config.method}("${config.path}", r.${config.funcName})\n`
        }
        out += line + '\n'
    }

    return out
}