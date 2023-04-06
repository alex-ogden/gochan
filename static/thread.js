
function showQuotedPost() {
    var quote = this;
    var parentParagraphId = quote.parentNode.id;
    var quotePreviewId = parentParagraphId.replace("-text", "-quotepreview");
    var quotePreview = document.getElementById(quotePreviewId);

    if (quote.className.includes("expanded")) {
      quotePreview.style.display = "none";
      quote.classList.remove("expanded");
    } else {
      var quotedPostNum = quote.innerHTML.replace("&gt;&gt;", "");
      var quotedText = document.getElementById(quotedPostNum+"-text").innerHTML;

      quotePreview.innerHTML = quotedText;
      quotePreview.style.display = "inline-block";
      quotePreview.style.border = "1px solid black";
      quotePreview.style.backgroundColor = "#363636"
      quotePreview.style.boxShadow = "2px 2px 4px rgba(0, 0, 0, 0.5)";
      quotePreview.style.padding = "30px";

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