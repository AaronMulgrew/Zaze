

async function hideTabPane()
{
  var i, items;
  items = document.getElementsByClassName("tab-pane");
  for (i = 0; i < items.length; i++)
  {
    items[i].style.display = 'none';
  }
}

// Operations when the web page is loaded.
async function onLoad(callback)
{
  await hideTabPane();
  document.getElementById("statusNotAuth").style.display = 'block';
  //document.getElementById("statusAuth").style.display = 'none';
  // Initiatlize CognitoAuth object
  auth = initCognitoSDK();
  document.getElementById("signInButton").addEventListener("click", function() {
    userButton(auth);
  });
  var curUrl = window.location.href;
  auth.parseCognitoWebResponse(curUrl);
  if (callback != null) {
    callback();    
  }
}

// Operation when tab is closed.
function closeTab(tabName) {
  var tab = document.getElementById(tabName);
  tab.style.display = 'none';
}

// Operation when tab is opened.
function openTab(tabName) {
  var tab = document.getElementById(tabName);
  tab.style.display = 'block';
}

// Perform user operations.
function userButton(auth) {
  var state = document.getElementById('signInButton').innerHTML;
  if (state === "Sign Out") {
    document.getElementById("signInButton").innerHTML = "Sign In";
    auth.signOut();
    showSignedOut();
  } else {
    auth.getSession();
  }
}

// Operations when signed in.
function showSignedIn(session)
{
  document.getElementById("statusNotAuth").style.display = 'none';
  //document.getElementById("statusAuth").style.display = 'block';
  document.getElementById("signInButton").innerHTML = "Sign Out";
  openTab("userdetails");
}

// Operations when signed out.
function showSignedOut()
{
  document.getElementById("statusNotAuth").style.display = 'block';
  //document.getElementById("statusAuth").style.display = 'none';
  closeTab("userdetails");
}

// Initialize a cognito auth object.
function initCognitoSDK() {
  var authData = {
    ClientId : '4psm863oa7vstodkvuaa7e6675', // Your client id here
    AppWebDomain : 'signin.zaze.io', // Exclude the "https://" part.
    TokenScopesArray : ['openid'], // like ['openid','email','phone']...
    RedirectUriSignIn : 'https://dev.zaze.io/',
    RedirectUriSignOut : 'https://dev.zaze.io/',
    //IdentityProvider : '<TODO: your identity provider you want to specify here>',
    //UserPoolId : '<TODO: your user pool id here>',
    //           AdvancedSecurityDataCollectionFlag : <TODO: boolean value indicating whether you want to enable advanced security data collection>
  };
  var auth = new AmazonCognitoIdentity.CognitoAuth(authData);
  // You can also set state parameter
  // auth.setState(<state parameter>);
  auth.userhandler = {
  // 	//onSuccess: <TODO: your onSuccess callback here>,
  // 	//onFailure: <TODO: your onFailure callback here>
      onSuccess: function(result) {
      //alert("Sign in success");
      showSignedIn(result);
    },
    onFailure: function(err) {
      alert("Error!" + err);
    }
  };
  // // The default response_type is "token", uncomment the next line will make it be "code".
  //auth.useCodeGrantFlow();
  return auth;
}
