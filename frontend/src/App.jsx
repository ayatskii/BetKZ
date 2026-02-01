import React from "react";

import Navbar from "./components/layout/Navbar";
import Register from "./features/auth/Register";
import EventList from "./features/betting/EventList";
import BetSlip from "./features/betting/BetSlip";
import Dashboard from "./features/dashboard/Dashboard";

export default function App() {
  return (
    <div>
      <Navbar />
      <div className="container">
        <Register />
        <EventList />
        <BetSlip />
        <Dashboard />
      </div>
    </div>
  );
}
