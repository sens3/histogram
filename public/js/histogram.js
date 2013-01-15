
function noopHandler(evt) {
  evt.stopPropagation();
  evt.preventDefault();
}

function displayImage(src) {
  var img = document.getElementById("preview"),
      histogram = document.getElementById("histogram");
  img.src = src;
  img.onload = function() {
    histogram.style.width = img.width;
    histogram.style.display = 'block';
  };
}

function handleReaderLoad(evt) {
  displayImage(evt.target.result);
}

function handleProgress(event) {
  if (event.lengthComputable) {
    var complete = (event.loaded / event.total * 100 | 0);
    progress.value = progress.innerHTML = complete;
  }
}

function getHistogramForForm(formData) {
	
	// now post a new XHR request
  var xhr = new XMLHttpRequest();
  xhr.open('POST', '/histogram');
  
  xhr.onloadstart = function () {
    progress.style.display = 'block';
    document.getElementById('histogram').style.display = 'none';
  }
  
  xhr.onload = function () {
    progress.style.display = 'none';
    if (xhr.status === 200) {
      progress.value = progress.innerHTML = 100;
      document.getElementById('histogram').innerHTML = xhr.response;
    } else {
      alert(xhr.response);
    }
  };
  
  xhr.upload.onprogress = handleProgress;

  xhr.send(formData)
}

function handleFile(file) {
  showFile(file)
  var formData = new FormData();
	formData.append('image_data', file);
  getHistogramForForm(formData);
}

function showFile(file) {
  var reader = new FileReader();

	// init the reader event handlers
	reader.onload = handleReaderLoad;

	// begin the read operation
	reader.readAsDataURL(file);
}

function drop(evt) {
	evt.stopPropagation();
	evt.preventDefault();

	var files = evt.dataTransfer.files;;

	// Only call the handler if 1 or more files was dropped.
	if (files.length > 0)
		handleFile(files[0]);
}

Histogram = {}

Histogram.exampleFileNames = ["chelsea-market.jpg", "nyc_1920.jpg", "pacman-ghost.jpg"];

Histogram.showExample = function() {
	var fileName = this.exampleFileNames[Math.floor(Math.random() * this.exampleFileNames.length)],
      formData = new FormData();
	
	displayImage("/images/" + fileName);
  
	formData.append('example_image_file', fileName);
  getHistogramForForm(formData);
};

Histogram.initEventHandlers = function() {
  var dropbox = document.getElementById("dropbox")

	// init event handlers
	dropbox.addEventListener("dragenter", noopHandler, false);
	dropbox.addEventListener("dragexit", noopHandler, false);
	dropbox.addEventListener("dragover", noopHandler, false);
	dropbox.addEventListener("drop", drop, false);
};

Histogram.init = function() {
  this.initEventHandlers();
  this.showExample();
};
