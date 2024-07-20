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
function config({ tableName }) {
    // convert tableName from snake case to camelCase
    const funcName = tableName.split('_').map((part) => {
        return part.charAt(0).toUpperCase() + part.slice(1)
    }).join('')

    // pluralize based on 

    return { funcName }
}

/**
 * @param {string} fileData
 * @param {ConfigObject} config
 * @returns {string} */
function addRoute(fileData, config) {
    const out = fileData


    tpl = `
-- name: ${config.funcName} :one
INSERT INTO import_logs (
    sync_started_at)
VALUES (
    ?)
RETURNING
    *;
`

    return out
}
