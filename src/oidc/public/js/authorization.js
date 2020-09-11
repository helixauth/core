

function signIn(e) {
  e.preventDefault()
  console.log("sign in");

  const emailInput = document.getElementById("email");
  const passwordInput = document.getElementById("password");

  fetch(`/authenticate?${QUERY}`, {
    method: "POST",
    body: JSON.stringify({ 
      email: emailInput ? emailInput.value : null,
      password: passwordInput ? passwordInput.value : null,
    }),
    headers: {
      "Content-Type": "application/json"
    },
  })
}


function signUp(e) {
  e.preventDefault()
  console.log("sign up");

  const emailInput = document.getElementById("email");
  const passwordInput = document.getElementById("password");

  fetch(`/authenticate?${QUERY}`, {
    method: "POST",
    body: JSON.stringify({ 
      email: emailInput ? emailInput.value : null,
      password: passwordInput ? passwordInput.value : null,
    }),
    headers: {
      "Content-Type": "application/json"
    },
  })
}
