function showData(data) {
  if (!data) {
    console.error("Data is empty");
    return;
  }

  if (!data.Followers) {
    console.error("No followers data found");
    return;
  }

  const followers = document.getElementById("followers");
  followers.innerHTML = data.Followers.total;
}

const collector = new TwitchCollector();
collector.registerHook("update", showData);
collector.start();
