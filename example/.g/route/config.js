/**
 * @param {Object} options
 * @param {string} options.method
 * @param {string} options.path
 * @returns {Object}
 */
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
