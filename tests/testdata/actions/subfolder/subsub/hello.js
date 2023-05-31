function main(args) {
    const name = args.name || "world";
    return { body: "Hello " + name + "!" };
}