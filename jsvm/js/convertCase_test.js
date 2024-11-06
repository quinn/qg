import { assertEquals } from "jsr:@std/assert";

Deno.test("kebab to camel", () => {
    eval(
        Deno.readTextFileSync("./convertCase.js")+
        'assertEquals(convertCase("pascal", "hello-world"), "HelloWorld");'+
        'assertEquals(convertCase("camel", "hello-world"), "helloWorld");'+
        'assertEquals(convertCase("snake", "hello-world"), "hello_world");'+
        'assertEquals(convertCase("kebab", "HelloWoRld"), "hello-wo-rld");'
    );
});
