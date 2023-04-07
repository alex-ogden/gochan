function showQuotedPost() {
    var quote = this;
    var parentParagraphId = quote.parentNode.id;
    var quotePreviewId = parentParagraphId.replace("-text", "-quotepreview");
    var quotePreviewImageId = parentParagraphId.replace("-text", "-quoteimage");
    var quotePreview = document.getElementById(quotePreviewId);
    var quotePreviewImage = document.getElementById(quotePreviewImageId);
    var quotedPostNum = quote.innerHTML.replace("&gt;&gt;", "");
    var quotedText = document.getElementById(quotedPostNum+"-text").innerHTML;
    // Images are the only element that have the post number as the ID so far
    var quotedImage = document.getElementById(quotedPostNum);
    var previewDivId = parentParagraphId.replace("-text", "-post-preview");
    var previewDiv = document.getElementById(previewDivId);

    if (quote.className.includes("expanded")) {
      previewDiv.removeAttribute("style");
      previewDiv.style.display = "none";
      quotePreview.removeAttribute("style");
      quote.classList.remove("expanded");
      // Some quoted elements might not have images
      if (quotedImage != undefined) {
        // Make the image disappear
        quotePreviewImage.src = "//:0";
        quotePreviewImage.removeAttribute("style");
      }
    } else {
      previewDiv.removeAttribute("style");
      previewDiv.style.border = "1px solid black";
      previewDiv.style.backgroundColor = "#363636"
      previewDiv.style.boxShadow = "2px 2px 4px rgba(0, 0, 0, 0.5)";

      // Some quoted elements might not have images
      if (quotedImage != undefined) {
        quotePreviewImage.removeAttribute("style");
        quotePreviewImage.src = quotedImage.src;
        quotePreviewImage.style.padding = "20px";
        quotePreviewImage.style.display = "inline";
        quotePreviewImage.style.boxShadow = "none";
      }

      quotePreview.innerHTML = quotedText;
      quotePreview.style.display = "inline-block";
      quotePreview.style.padding = "20px";

      // Add expanded class
      quote.classList.add("expanded");
    }
}

var quotes = document.getElementsByClassName("quotelink");
for (var i = 0; i < quotes.length; i++) {
    /*
        Remove href attribute to stop browser from counting
        quote links as actual links and reloading the page
    */
    quotes[i].removeAttribute("href");
    quotes[i].addEventListener("click", showQuotedPost);
}