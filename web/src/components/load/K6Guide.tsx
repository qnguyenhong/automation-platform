import { BookOpen, HelpCircle } from 'lucide-react'

export default function K6Guide() {
  return (
    <div className="glass-panel text-slate-800 rounded-xl p-6 border border-slate-200 shadow-sm space-y-6">
      <div className="flex items-center gap-3 border-b border-slate-200 pb-4">
        <BookOpen className="w-5 h-5 text-indigo-600" />
        <div>
          <h3 className="font-bold text-base text-slate-800">K6 Configuration Guide</h3>
          <p className="text-xs text-slate-500">Learn how to translate standard k6 scripts into our high-performance concurrent engine.</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Left Side: K6 JavaScript Script */}
        <div className="space-y-2">
          <span className="text-xs font-bold uppercase tracking-wider text-indigo-600">1. Your k6 Script (options)</span>
          <pre className="p-4 bg-slate-50 rounded-lg text-xs font-mono overflow-x-auto text-slate-700 border border-slate-200 h-64 shadow-inner">
{`import http from 'k6/http';

export const options = {
  vus: 50,
  duration: '30s',
  stages: [
    { duration: '10s', target: 50 }, // Ramp up
  ],
  thresholds: {
    http_req_failed: ['rate<0.01'],  // error rate SLO
    http_req_duration: ['p(95)<200'] // latency SLO
  }
};

export default function () {
  http.get('https://api.example.com/users');
}`}
          </pre>
        </div>

        {/* Right Side: Our Configuration */}
        <div className="space-y-2">
          <span className="text-xs font-bold uppercase tracking-wider text-emerald-600">2. Equivalent Native Config</span>
          <pre className="p-4 bg-slate-50 rounded-lg text-xs font-mono overflow-x-auto text-slate-700 border border-slate-200 h-64 shadow-inner">
{`{
  "method": "GET",
  "url": "https://api.example.com/users",
  "load": {
    "vus": 50,
    "duration": "30s",
    "ramp_up": "10s",
    "ramp_down": "0s",
    "rate_limit_rps": 0
  },
  "assertions": [
    {
      "type": "error_rate",
      "expected": 0.01
    },
    {
      "type": "p95_response_time",
      "expected": 200
    }
  ]
}`}
          </pre>
        </div>
      </div>

      {/* Mapping Details */}
      <div className="space-y-3 pt-2">
        <h4 className="text-xs font-bold uppercase tracking-wider text-slate-700 flex items-center gap-1">
          <HelpCircle className="w-4 h-4 text-slate-500" />
          Translation Key Matrix
        </h4>
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs border-collapse">
            <thead>
              <tr className="border-b border-slate-200 text-slate-550 font-semibold">
                <th className="py-2 pr-4">k6 Concept</th>
                <th className="py-2 px-4">Our Platform parameter</th>
                <th className="py-2 pl-4">Description</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100 text-slate-700 font-mono">
              <tr>
                <td className="py-2.5 pr-4 font-bold text-indigo-600">options.vus</td>
                <td className="py-2.5 px-4 font-bold text-emerald-600">vus</td>
                <td className="py-2.5 pl-4 font-sans text-slate-650">Number of concurrent Virtual User routines to spawn.</td>
              </tr>
              <tr>
                <td className="py-2.5 pr-4 font-bold text-indigo-600">options.duration</td>
                <td className="py-2.5 px-4 font-bold text-emerald-600">duration</td>
                <td className="py-2.5 pl-4 font-sans text-slate-650">Total duration of the load test (e.g. 10s, 5m, 1h).</td>
              </tr>
              <tr>
                <td className="py-2.5 pr-4 font-bold text-indigo-600">stages (ramp-up)</td>
                <td className="py-2.5 px-4 font-bold text-emerald-600">ramp_up</td>
                <td className="py-2.5 pl-4 font-sans text-slate-650">Duration over which the active VU count increases linearly.</td>
              </tr>
              <tr>
                <td className="py-2.5 pr-4 font-bold text-indigo-600">stages (ramp-down)</td>
                <td className="py-2.5 px-4 font-bold text-emerald-600">ramp_down</td>
                <td className="py-2.5 pl-4 font-sans text-slate-650">Duration over which the active VU count decreases linearly to 0.</td>
              </tr>
              <tr>
                <td className="py-2.5 pr-4 font-bold text-indigo-600">options.rps / rate</td>
                <td className="py-2.5 px-4 font-bold text-emerald-600">rate_limit_rps</td>
                <td className="py-2.5 pl-4 font-sans text-slate-650">Total request rate limit (RPS) distributed across all VUs.</td>
              </tr>
              <tr>
                <td className="py-2.5 pr-4 font-bold text-indigo-600">thresholds</td>
                <td className="py-2.5 px-4 font-bold text-emerald-600">assertions</td>
                <td className="py-2.5 pl-4 font-sans text-slate-650">SLO validation (e.g. error rate, p95 response time, minimum throughput).</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
