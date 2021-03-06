var HTMLContents = "";
var PostName = "";



function completeAndRedirect(auth)
{
  var formElements = document.getElementById("designForm").elements;
  var jsonData = {};

  for (var i = 0; i < formElements.length; i++)
  {
    var elementName = formElements[i].name;
    jsonData[elementName] = formElements[i].value;
  }
  var myEditor = document.querySelector('#editor-container')
  var content = myEditor.children[0].innerHTML
  jsonData["content"] = content

  var jwtToken = auth.getSignInUserSession().idToken.jwtToken
  var xhr = new XMLHttpRequest();
  var url = "https://5o3kiu1m91.execute-api.eu-west-2.amazonaws.com/dev/generate_static_site";
  xhr.open("POST", url, true);
  xhr.setRequestHeader("Authorization", "Bearer " + jwtToken);
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.onreadystatechange = function () {
    if (xhr.readyState === 4 && xhr.status === 200) {
      //console.log(xhr.responseText);
	    //document.body.innerHTML = xhr.responseText
      closeTab("userdetails");
      openTab("preview");
      document.getElementById('title').innerHTML = "Preview Designs";
      HTMLContents = xhr.responseText;
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

function displaySuccess(link)
{
  closeTab("preview");
  document.getElementById("success_link").innerHTML = link;
  document.getElementById("success_link").href = link;
  user_index_link = "https://users.zaze.io/" + auth.username + "/index.html";
  document.getElementById("index_link").innerHTML = user_index_link;
  document.getElementById("index_link").href = user_index_link;
  openTab("success");
}

function getEditedPageContent(pageName)
{
  PostName = pageName;
  var jwtToken = auth.getSignInUserSession().idToken.jwtToken
  var xhr = new XMLHttpRequest();
  var url = "https://5o3kiu1m91.execute-api.eu-west-2.amazonaws.com/dev/get_one_post";
  xhr.open("POST", url, true);
  xhr.setRequestHeader("Authorization", "Bearer " + jwtToken);
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.onreadystatechange = function () {
    if (xhr.readyState === 4 && xhr.status === 200) {
      //alert(xhr.responseText);
      var contents = JSON.parse(xhr.responseText);
      document.getElementById("header").value = contents[0];
      quill.clipboard.dangerouslyPasteHTML(contents[1]);
      //document.getElementById("editor-container").innerHTML = contents[1];
      document.getElementById("background-color").value = contents[2];
      document.getElementById("font-color").value = contents[3];
    }
  };
  var jsonData = {'postname':pageName};
  var data = JSON.stringify(jsonData);
  xhr.send(data);
}


function designOnLoad()
{
  onLoad(function(){
    try {
      edit_page_url = location.hash.split('page=')[1].split('&')[0];
    }
    catch {
      edit_page_url = null;
    }
    if(edit_page_url != null)
    {
      getEditedPageContent(edit_page_url);
    }
  });
}


function proceedToGeneration()
{
  var title = document.getElementById("designForm").elements["header"].value;
  var jsonData = {};
  jsonData["HTMLContents"] = HTMLContents;
  jsonData["title"] = title;
  if (PostName != "")
  {
    jsonData['Edit'] = true;
    jsonData['UniqueID'] = PostName;
  }

  var jwtToken = auth.getSignInUserSession().idToken.jwtToken
  var xhr = new XMLHttpRequest();
  var url = "https://5o3kiu1m91.execute-api.eu-west-2.amazonaws.com/dev/post_to_s3";
  xhr.open("POST", url, true);
  xhr.setRequestHeader("Authorization", "Bearer " + jwtToken);
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.onreadystatechange = function () {
    if (xhr.readyState === 4 && xhr.status === 200) {
      GenerateIndex(xhr.responseText);
    }
  };
  delete jsonData[""];
  console.log(jsonData)
  var data = JSON.stringify(jsonData);
  xhr.send(data);
}




function GenerateIndex(success_link)
{
  var xhr = new XMLHttpRequest();
  var url = "https://5o3kiu1m91.execute-api.eu-west-2.amazonaws.com/dev/generate_index";
  var jwtToken = auth.getSignInUserSession().idToken.jwtToken
  xhr.open("POST", url, true);
  xhr.setRequestHeader("Authorization", "Bearer " + jwtToken);
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.onreadystatechange = function () {
    if (xhr.readyState === 4 && xhr.status === 200) {
      console.log(xhr.responseText);
      displaySuccess(success_link);
    }
  };
  xhr.send();
}
