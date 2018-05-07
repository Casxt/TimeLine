function urlencodeFormData(fd) {
    var params = new URLSearchParams();
    for (var pair of fd.entries()) {
        typeof pair[1] == 'string' && params.append(pair[0], pair[1]);
    }
    return params.toString();
}