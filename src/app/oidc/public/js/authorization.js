const signInForm = document.getElementById("sign-in-form");
if (signInForm) {
  signInForm.addEventListener("submit", function(event) {
    event.preventDefault();
    console.log("sign in");
    // TODO
  });
}

const signUpForm = document.getElementById("sign-up-form");
if (signUpForm) {
  signUpForm.addEventListener("submit", function(event) {
    event.preventDefault();
    console.log("sign up");
    // TODO
  });
}
