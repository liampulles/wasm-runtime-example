function handleRequest(env, req) {
    var path = req.URL.Path;
    env.Printf("Calling %s...\n", path);
}