import { goto } from "$app/navigation";

export function logout() {
    localStorage.removeItem('token');

    fetch('http://localhost:3030/logout', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    })
    .then(response => response.json())
    .then(data => {
      console.log(data.message); 
        goto("/auth/login")
    })
    .catch(error => console.error('Error:', error));
  }
  