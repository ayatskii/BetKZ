import React, { useState } from "react";
import { createUser } from "../../services/api";

export default function Register() {
  const [email, setEmail] = useState("");
  const [balance, setBalance] = useState("");
  const [user, setUser] = useState(null);

  const submit = async () => {
    const result = await createUser({
      email,
      balance: Number(balance),
    });
    setUser(result);
  };

  return (
    <section>
      <h2>Register User</h2>
      <input placeholder="Email" onChange={e => setEmail(e.target.value)} />
      <input
        type="number"
        placeholder="Balance"
        onChange={e => setBalance(e.target.value)}
      />
      <button onClick={submit}>Create</button>
      {user && <pre>{JSON.stringify(user, null, 2)}</pre>}
    </section>
  );
}
