
function noopHandler(evt) {
  evt.stopPropagation();
  evt.preventDefault();
}

function handleReaderLoad(evt) {
  var img = document.getElementById("preview");
  img.src = evt.target.result;
  setTimeout(function(){
    var histogram = document.getElementById("histogram");
    histogram.style.width = img.width;
    
  }, 5);
}

function uploadFile(file) {
  var formData = new FormData();
	formData.append('image_file', file);
	
	// now post a new XHR request
  var xhr = new XMLHttpRequest();
  xhr.open('POST', '/histogram');
  
  xhr.onloadstart = function () {
    progress.style.display = 'block';
  }
  
  xhr.onload = function () {
    if (xhr.status === 200) {
      progress.value = progress.innerHTML = 100;
      progress.style.display = 'none';
      document.getElementById('histogram').innerHTML = xhr.response;
    } else {
      alert("Error uploading image");
    }
  };
  
  xhr.upload.onprogress = function (event) {
    if (event.lengthComputable) {
      var complete = (event.loaded / event.total * 100 | 0);
      progress.value = progress.innerHTML = complete;
    }
  };

  xhr.send(formData)
}

function handleFile(file) {

	var reader = new FileReader();

	// init the reader event handlers
	reader.onload = handleReaderLoad;

	// begin the read operation
	reader.readAsDataURL(file);
	
  uploadFile(file);
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

Histogram.init = function() {

	var dropbox = document.getElementById("dropbox")

	// init event handlers
	dropbox.addEventListener("dragenter", noopHandler, false);
	dropbox.addEventListener("dragexit", noopHandler, false);
	dropbox.addEventListener("dragover", noopHandler, false);
	dropbox.addEventListener("drop", drop, false);
	
}