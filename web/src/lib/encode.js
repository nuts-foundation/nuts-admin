// encodeURIPath is like encodeURI, but does not encode :
// Reason for this, is that Verifiable Credential IDs (VC IDs) are URIs that contain colons, which don't need to be encoded in URI paths,
// but are encoded by encodeURI. This breaks API usage in Go, because Echo (or Go's HTTP package) does not decode the colons in the path, leading to 404 errors when trying to access VC IDs
export function encodeURIPath(path) {
    // for-loop to iterate over each character in the path
    let encodedPath = '';
    for (let i = 0; i < path.length; i++) {
        const char = path[i];
        // encode all characters except for alphanumeric characters, -, _, ., ~, and :
        if (/^[a-zA-Z0-9\-_.~:]$/.test(char)) {
            encodedPath += char;
        } else {
            encodedPath += encodeURIComponent(char);
        }
    }
    return encodedPath;
}