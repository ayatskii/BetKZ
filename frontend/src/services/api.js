const API_URL = "http://localhost:8080";


export async function createUser(data) {
  const res = await fetch(`${API_URL}/users`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  return res.json();
}

export async function getEvents() {
  const res = await fetch(`${API_URL}/events`);
  return res.json();
}

export async function placeBet(data) {
  const res = await fetch(`${API_URL}/bets`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error("Bet failed");
  return res.json();
}
