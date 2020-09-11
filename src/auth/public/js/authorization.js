

function signIn(e) {
  e.preventDefault()
  console.log("sign in");
  console.log("help");

  const emailInput = document.getElementById("email");
  const passwordInput = document.getElementById("password");
  const errorParagraph = document.getElementById("error");

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
  .then(res => res.json())
  .then(res => {
    console.log(res)
    if (res.error) {
      errorParagraph.innerText = res.error
    }
  })
  .catch(err => {
    console.log(err)
    errorParagraph.innerText = err
  })
}


function signUp(e) {
  e.preventDefault()
  console.log("sign up");

  const emailInput = document.getElementById("email");
  const passwordInput = document.getElementById("password");
  const errorParagraph = document.getElementById("error");

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
  .then(res => res.json())
  .then(res => {
    console.log(res)
    if (res.error) {
      errorParagraph.innerText = res.error
    }
  })  
  .catch(err => {
    console.log(err)
    errorParagraph.innerText = err
  })
}
