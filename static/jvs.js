function show(artistName) {
    window.location.href = '/details.html?artist=' + encodeURIComponent(artistName);
}
function home() {
    window.location.href = '/';
}