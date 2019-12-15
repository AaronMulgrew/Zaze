insert = function insert(main_string, ins_string, pos) {
   if(typeof(pos) === "undefined") {
    pos = 0;
  }
   if(typeof(ins_string) === "undefined") {
    ins_string = '';
  }
   return main_string.slice(0, pos) + ins_string + main_string.slice(pos);
}


function goto_edit_post(calling_element)
{
  var url = window.location.href.toString();
  var urlIndexInsert = url.indexOf("zaze.io/") + 8;
  url = insert(url, "designer.html", urlIndexInsert);


  if (calling_element != undefined)
  {
    calling_elem_id = calling_element.id;
    url = url += '&page=' + calling_elem_id;
  }
  window.location.href = url;

}


function generate_card(uniqueID, obj_name)
{
  cardsdiv = document.getElementById("cardsdiv");
  base_div = document.createElement("div");
  base_div.setAttribute("class", "card text-white bg-dark mb-3 mx-auto");
  card_body = document.createElement("div");
  card_body.setAttribute("class", "card-body");
  title = document.createElement("h5");
  title.innerHTML = obj_name;
  edit_button = document.createElement("a");
  edit_button.id = uniqueID;
  edit_button.setAttribute("onClick", "goto_edit_post(this)");
  edit_button.setAttribute("class", "btn btn-primary");
  edit_button.innerHTML = "edit post";
  delete_button = document.createElement("button");
  delete_button.setAttribute("class", "btn btn-danger");
  delete_button.id = uniqueID;
  delete_button.setAttribute("onclick", "showModal(this);");

  delete_button.innerHTML = "delete post";
  card_body.appendChild(title);
  card_body.appendChild(edit_button);
  card_body.appendChild(delete_button);
  base_div.appendChild(card_body);
  cardsdiv.appendChild(base_div);
}

function delete_post(pageName)
{
  var jwtToken = auth.getSignInUserSession().idToken.jwtToken
  var xhr = new XMLHttpRequest();
  var url = "https://5o3kiu1m91.execute-api.eu-west-2.amazonaws.com/dev/delete_post";
  xhr.open("POST", url, true);
  xhr.setRequestHeader("Authorization", "Bearer " + jwtToken);
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.onreadystatechange = function () {
    if (xhr.readyState === 4 && xhr.status === 200) {
      $('#exampleModal').modal('hide');
      render_all_posts();
    }
  };
  var jsonData = {"PageName": pageName};
  var data = JSON.stringify(jsonData);
  xhr.send(data);
}

function showModal(calling_element)
{
  calling_elem_id = calling_element.id;
  document.getElementById("deleteButton").setAttribute("onclick", "delete_post(\"" + calling_elem_id + "\")");
  $('#exampleModal').modal('show')

}

function render_all_posts()
{
  var jwtToken = auth.getSignInUserSession().idToken.jwtToken;
  var xhr = new XMLHttpRequest();
  var url = "https://5o3kiu1m91.execute-api.eu-west-2.amazonaws.com/dev/get_current_posts";
  xhr.open("POST", url, true);
  xhr.setRequestHeader("Authorization", "Bearer " + jwtToken);
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.onreadystatechange = function () {
    if (xhr.readyState === 4 && xhr.status === 200) {

      // clear out the user details children
      cardsdiv = document.getElementById("cardsdiv");
      while (cardsdiv.firstChild) {
        cardsdiv.removeChild(cardsdiv.firstChild);
      }

      var bucket_items = JSON.parse(xhr.responseText);

      for ( var i = 0; i < bucket_items.length; i++)
      {
        var post_title = bucket_items[i].PostTitle;
        var uniqueID = bucket_items[i].UniqueID;
        if(post_title != "")
        {
          if(uniqueID != "")
          {
            generate_card(uniqueID, post_title);
          }
        }
      }
    }
    else if (xhr.readyState === 4 && xhr.status != 200) {
      alert("Credentials timeout");
      auth.signOut();
      showSignedOut();
    }
  };
  xhr.send();
}

function onload_all_posts()
{
  onLoad(function(){
    if(auth.getCurrentUser() != null) {
      render_all_posts();
    }
  });
}
