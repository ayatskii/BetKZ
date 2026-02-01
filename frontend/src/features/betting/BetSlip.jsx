import React, { useState } from "react";
import { placeBet } from "../../services/api";

export default function BetSlip() {
  const [form, setForm] = useState({});
  const [bet, setBet] = useState(null);

  const submit = async () => {
    const result = await placeBet({
      user_id: Number(form.user_id),
      event_id: Number(form.event_id),
      amount: Number(form.amount),
    });
    setBet(result);
  };

  return (
    <section>
      <h2>Bet Slip</h2>
      <input onChange={e => setForm({ ...form, user_id: e.target.value })} placeholder="User ID" />
      <input onChange={e => setForm({ ...form, event_id: e.target.value })} placeholder="Event ID" />
      <input onChange={e => setForm({ ...form, amount: e.target.value })} placeholder="Amount" />
      <button onClick={submit}>Place Bet</button>
      {bet && <pre>{JSON.stringify(bet, null, 2)}</pre>}
    </section>
  );
}
