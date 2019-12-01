function completeAndRedirect(auth)
{
  var formElements = document.getElementById("designForm").elements;
  var jsonData = {};

  for (var i = 0; i < formElements.length; i++)
  {
    var elementName = formElements[i].name;
    jsonData[elementName] = formElements[i].value;
  }

  var jwtToken = auth.getSignInUserSession().idToken.jwtToken
  var xhr = new XMLHttpRequest();
  var url = "https://5o3kiu1m91.execute-api.eu-west-2.amazonaws.com/dev/generate_static_site";
  xhr.open("POST", url, true);
  xhr.setRequestHeader("Authorization", "Bearer " + jwtToken);
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.onreadystatechange = function () {
    if (xhr.readyState === 4 && xhr.status === 200) {
      console.log(xhr.responseText);
	    //document.body.innerHTML = xhr.responseText
      closeTab("userdetails");
      openTab("preview");
      document.getElementById('title').innerHTML = "Preview Designs"
      document.getElementById('preview_site').src = "data:text/html;charset=utf-8," + escape(xhr.responseText);
    }
  };
  delete jsonData[""];
  console.log(jsonData)
  var data = JSON.stringify(jsonData);
  xhr.send(data);
}

function revertToDesignTab()
{
  openTab("userdetails");
  closeTab("preview");
}

function proceedToGeneration()
{
  var formElements = document.getElementById("designForm").elements;
  var jsonData = {};

  for (var i = 0; i < formElements.length; i++)
  {
    var elementName = formElements[i].name;
    jsonData[elementName] = formElements[i].value;
  }

  var jwtToken = auth.getSignInUserSession().idToken.jwtToken
  var xhr = new XMLHttpRequest();
  var url = "https://5o3kiu1m91.execute-api.eu-west-2.amazonaws.com/dev/post_to_s3";
  xhr.open("POST", url, true);
  xhr.setRequestHeader("Authorization", "Bearer " + jwtToken);
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.onreadystatechange = function () {
    if (xhr.readyState === 4 && xhr.status === 200) {
      console.log(xhr.responseText);
	    //document.body.innerHTML = xhr.responseText
      closeTab("userdetails");
      openTab("preview");
      document.getElementById('title').innerHTML = "Preview Designs"
      document.getElementById('preview_site').src = "data:text/html;charset=utf-8," + escape(xhr.responseText);
    }
  };
  delete jsonData[""];
  console.log(jsonData)
  var data = JSON.stringify(jsonData);
  xhr.send(data);
}
