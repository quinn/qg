function detectCase(input) {
    if (input.includes('-')) return 'kebab';
    if (input.includes('_')) return 'snake';
    if (/^[A-Z]/.test(input) && !/[a-z]/.test(input)) return 'pascal';
    if (/^[a-z]/.test(input) && /[A-Z]/.test(input)) return 'camel';
    if (/^[A-Z]/.test(input) && /[a-z]/.test(input)) return 'pascal';
    return 'unknown';
}

function toWords(input, inputCase) {
    switch (inputCase) {
        case 'kebab':
            return input.split('-');
        case 'snake':
            return input.split('_');
        case 'camel':
            return input.replace(/([a-z])([A-Z])/g, '$1 $2').split(' ');
        case 'pascal':
            return input.replace(/([A-Z][a-z]*)/g, ' $1').trim().split(' ');
        default:
            return [input];
    }
}

function toKebabCase(words) {
    return words.map(word => word.toLowerCase()).join('-');
}

function toSnakeCase(words) {
    return words.map(word => word.toLowerCase()).join('_');
}

function toCamelCase(words) {
    return words.map((word, index) => {
        if (index === 0) return word.toLowerCase();
        return word.charAt(0).toUpperCase() + word.slice(1).toLowerCase();
    }).join('');
}

function toPascalCase(words) {
    return words.map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase()).join('');
}

function convertCase(targetCase, input) {
    const inputCase = detectCase(input);
    const words = toWords(input, inputCase);

    switch (targetCase) {
        case 'kebab':
            return toKebabCase(words);
        case 'snake':
            return toSnakeCase(words);
        case 'camel':
            return toCamelCase(words);
        case 'pascal':
            return toPascalCase(words);
        default:
            return input;
    }
}

// Examples:
// console.log(convertCase('camel', 'this-is-kebab-case')); // thisIsKebabCase
// console.log(convertCase('pascal', 'this_is_snake_case')); // ThisIsSnakeCase
// console.log(convertCase('snake', 'thisIsCamelCase')); // this_is_camel_case
// console.log(convertCase('kebab', 'ThisIsPascalCase')); // this-is-pascal-case
