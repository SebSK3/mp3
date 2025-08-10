const progressSize = 20;
let advancedSpan = document.getElementById("advanced");
let progress = document.getElementById("progress");
let list = document.getElementById("list");
let intervalId = null;

function toggleAdvanced() {
  if (advancedSpan.innerHTML != "") {
    advancedSpan.innerHTML = "";
    return;
  }
  let advanced = ' \
    <input type="text" placeholder="Artist" name="Artist"> \
    <input type="text" placeholder="Title" name="Title"> \
    <input type="text" placeholder="Album" name="Album"> \
    <input type="text" placeholder="Custom Image URL" name="ImgUrl"> \
    <input type="text" placeholder="Custom FFmpeg command" name="FfmpegCmd">';
  advancedSpan.innerHTML = advanced;
}

function checkProgress(id) {
  fetch(`/status/${id}`)
  .then(response => response.json())  
  .then(data => {
    const currentProgress = data.progress || 0;
      progress.innerHTML = "[" +
        "=".repeat(currentProgress / (100 / progressSize)) +
        " ".repeat((100-currentProgress) / (100 / progressSize)) +
        "]";
      if (currentProgress >= 100) {
        clearInterval(intervalId);
      }
  })
  .catch(err => {
    // TODO error handling
    clearInterval(intervalId);
  })
}

function startProgress(id) {
  console.log("Starting progress check on id: " + id);
  intervalId = setInterval(() => checkProgress(id), 1000);
}


function getList() {
  if (!document.URL.includes("list")) {
    return;
  }
  
}

// Everything below is vibe coded
document.getElementById("download-form").addEventListener("submit", async function (e) {
  e.preventDefault();

  const form = e.target;
  const formData = new FormData(form);
  const params = new URLSearchParams();

  for (const [key, value] of formData.entries()) {
    params.append(key, value);
  }
  fetch("/download", {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: params,
  })
  .then(response => {
    if (!response.ok) {
      throw new Error("Server returned an error");
    }
    return response.json();
  })
  .then(data => {
    startProgress(data.id);
  })
  .catch(error => {
    console.error("Error:", error);
    // TODO: here handling error messages
  });
  
  });
getList();
