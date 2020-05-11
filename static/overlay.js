function showData(data) {
    if (typeof data === "undefined" || data === null) {
        console.error("Data is empty");
        return;
    }

    if (typeof data.Followers === "undefined" || data.Followers === null) {
        console.error("No followers data found");
        return;
    }

    let followers = document.getElementById("followers");
    followers.innerHTML = data.Followers.total;
}

let collector = Object.create(twitchCollector);
collector.registerHook("update", function(data) {
    showData(data);
});
collector.start();