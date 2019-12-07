
function generate_card(obj)
{
  userdetails = document.getElementById("userdetails");

  base_div = document.createElement("div");
  base_div.setAttribute("class", "card text-white bg-dark mb-3 mx-auto");
  card_body = document.createElement("div");
  card_body.setAttribute("class", "card-body");
  title = document.createElement("h5");
  title.innerHTML = obj;
  edit_button = document.createElement("a");
  edit_button.href = "zaze.io/testing123";
  edit_button.setAttribute("class", "btn btn-primary");
  edit_button.innerHTML = "edit post";
  delete_button = document.createElement("a");
  delete_button.href = "zaze.io/testing123";
  delete_button.setAttribute("class", "btn btn-danger");
  delete_button.innerHTML = "delete post";
  card_body.appendChild(title);
  card_body.appendChild(edit_button);
  card_body.appendChild(delete_button);
  base_div.appendChild(card_body);
  userdetails.appendChild(base_div);
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
      console.log(xhr.responseText);
      var bucket_items = JSON.parse(xhr.responseText);

      for ( var i = 0; i < bucket_items.length; i++)
      {
        var obj = bucket_items[i];

        if(obj != "")
        {
          generate_card(obj);
        }
      }
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
