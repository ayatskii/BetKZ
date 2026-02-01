import React, { useState } from "react";
import { getEvents } from "../../services/api";

export default function EventList() {
  const [events, setEvents] = useState([]);

  const load = async () => {
    const data = await getEvents();
    setEvents(data);
  };

  return (
    <section>
      <h2>Events</h2>
      <button onClick={load}>Load Events</button>
      <ul>
        {events.map(e => (
          <li key={e.id}>
            #{e.id} — {e.name} (odds: {e.odds})
          </li>
        ))}
      </ul>
    </section>
  );
}
