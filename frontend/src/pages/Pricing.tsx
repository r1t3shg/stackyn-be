import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import Logo from '@/components/Logo';

export default function Pricing() {
  const { user } = useAuth();
  const navigate = useNavigate();

  const handleSignInClick = (e: React.MouseEvent<HTMLAnchorElement>) => {
    e.preventDefault();
    if (user) {
      window.location.href = 'https://console.staging.stackyn.com/';
    } else {
      navigate('/login');
    }
  };

  return (
    <div className="min-h-screen bg-[var(--app-bg)]">
      {/* Header */}
      <header className="border-b border-[var(--border-subtle)] bg-[var(--app-bg)] sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center gap-8">
              {/* Logo */}
              <div className="flex-shrink-0">
                <Logo />
              </div>

              {/* Desktop Navigation Links - Left Aligned */}
              <nav className="hidden md:flex items-center space-x-6">
                <a
                  href="/#why-stackyn"
                  className="text-[var(--text-secondary)] hover:text-[var(--text-primary)] transition-colors font-medium"
                >
                  Why Stackyn?
                </a>
                <a
                  href="/#features"
                  className="text-[var(--text-secondary)] hover:text-[var(--text-primary)] transition-colors font-medium"
                >
                  Features
                </a>
                <Link
                  to="/pricing"
                  className="text-[var(--primary)] font-medium"
                >
                  Pricing
                </Link>
              </nav>
            </div>

            {/* Desktop CTA - Right Side */}
            <div className="hidden md:flex items-center gap-3">
              {user ? (
                <>
                  <a
                    href="https://console.staging.stackyn.com/"
                    className="bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-[var(--app-bg)] font-medium py-2 px-6 rounded-lg transition-colors"
                  >
                    Go to Console
                  </a>
                  <button
                    className="w-10 h-10 rounded-full bg-[var(--primary-muted)] flex items-center justify-center text-[var(--primary)] font-semibold hover:bg-[var(--elevated)] transition-colors"
                    aria-label="User menu"
                  >
                    {user.email?.charAt(0).toUpperCase() || 'U'}
                  </button>
                </>
              ) : (
                <a
                  href="/login"
                  onClick={handleSignInClick}
                  className="bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-[var(--app-bg)] font-medium py-2 px-6 rounded-lg transition-colors"
                >
                  Sign in
                </a>
              )}
            </div>

            {/* Mobile menu button */}
            <div className="md:hidden flex items-center">
              <button
                className="text-[var(--text-secondary)] hover:text-[var(--text-primary)] p-2"
                aria-label="Toggle menu"
              >
                <svg
                  className="h-6 w-6"
                  fill="none"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path d="M4 6h16M4 12h16M4 18h16" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Content */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        {/* Header Section */}
        <div className="text-center mb-16">
          <h1 className="text-4xl sm:text-5xl font-bold text-[var(--text-primary)] mb-4">Pricing</h1>
          <h2 className="text-2xl sm:text-3xl font-semibold text-[var(--text-primary)] mb-4">
            Simple pricing. No surprises.
          </h2>
          <p className="text-xl text-[var(--text-secondary)] max-w-3xl mx-auto">
            Pay only for the resources your apps actually use.
          </p>
        </div>

        {/* Plans Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-20">
          {/* Free Plan */}
          <div className="border border-[var(--border-subtle)] rounded-lg p-6 bg-[var(--surface)]">
            <h3 className="text-2xl font-bold text-[var(--text-primary)] mb-2">Free</h3>
            <p className="text-[var(--text-muted)] mb-6 text-sm">For trying Stackyn and small experiments.</p>
            <ul className="space-y-3 mb-8">
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">1 app</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">256 MB RAM</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">Shared CPU</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">HTTPS included</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">Manual deploys</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">Basic logs</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">Community support</span>
              </li>
            </ul>
            <div className="mb-6">
              <div className="text-3xl font-bold text-[var(--text-primary)] mb-2">$0</div>
            </div>
            <a
              href={user ? "https://console.staging.stackyn.com/" : "/login"}
              onClick={handleSignInClick}
              className="block w-full text-center bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-[var(--app-bg)] font-medium py-3 px-6 rounded-lg transition-colors"
            >
              Get Started
            </a>
          </div>

          {/* Pro Plan */}
          <div className="border-2 border-[var(--primary)] rounded-lg p-6 bg-[var(--elevated)] relative">
            <div className="absolute top-0 right-0 bg-[var(--primary)] text-[var(--app-bg)] px-3 py-1 rounded-bl-lg rounded-tr-lg text-xs font-medium">
              Popular
            </div>
            <h3 className="text-2xl font-bold text-[var(--text-primary)] mb-2">Pro</h3>
            <p className="text-[var(--text-muted)] mb-6 text-sm">For indie hackers and production apps.</p>
            <ul className="space-y-3 mb-8">
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">Up to 5 apps</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">512 MB RAM per app</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">Automatic deploys on Git push</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">HTTPS + custom domains</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">Real-time logs</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">App health status</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">Email support</span>
              </li>
            </ul>
            <div className="mb-6">
              <div className="text-3xl font-bold text-[var(--text-primary)] mb-2">$15 <span className="text-lg font-normal text-[var(--text-muted)]">/ month</span></div>
            </div>
            <a
              href={user ? "https://console.staging.stackyn.com/" : "/login"}
              onClick={handleSignInClick}
              className="block w-full text-center bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-[var(--app-bg)] font-medium py-3 px-6 rounded-lg transition-colors"
            >
              Start Pro
            </a>
          </div>

          {/* Team Plan */}
          <div className="border border-[var(--border-subtle)] rounded-lg p-6 bg-[var(--surface)]">
            <h3 className="text-2xl font-bold text-[var(--text-primary)] mb-2">Team</h3>
            <p className="text-[var(--text-muted)] mb-6 text-sm">For small teams and agencies.</p>
            <ul className="space-y-3 mb-8">
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">Up to 20 apps</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">1 GB RAM per app</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">Automatic deploys</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">HTTPS + custom domains</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">Centralized logs</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">Priority support</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">Team access</span>
              </li>
            </ul>
            <div className="mb-6">
              <div className="text-3xl font-bold text-[var(--text-primary)] mb-2">$49 <span className="text-lg font-normal text-[var(--text-muted)]">/ month</span></div>
            </div>
            <a
              href={user ? "https://console.staging.stackyn.com/" : "/login"}
              onClick={handleSignInClick}
              className="block w-full text-center bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-[var(--app-bg)] font-medium py-3 px-6 rounded-lg transition-colors"
            >
              Start Team
            </a>
          </div>

          {/* Custom Plan */}
          <div className="border border-[var(--border-subtle)] rounded-lg p-6 bg-[var(--surface)]">
            <h3 className="text-2xl font-bold text-[var(--text-primary)] mb-2">Custom</h3>
            <p className="text-[var(--text-muted)] mb-6 text-sm">For higher scale or special requirements.</p>
            <ul className="space-y-3 mb-8">
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">More apps</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">Higher RAM limits</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">Dedicated resources</span>
              </li>
              <li className="flex items-start text-sm">
                <svg className="w-5 h-5 text-[var(--success)] mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-[var(--text-secondary)]">Custom setup</span>
              </li>
            </ul>
            <div className="mb-6">
              <div className="text-3xl font-bold text-[var(--text-primary)] mb-2">Custom</div>
            </div>
            <a
              href="#contact"
              className="block w-full text-center bg-[var(--surface)] hover:bg-[var(--elevated)] text-[var(--text-primary)] font-medium py-3 px-6 rounded-lg transition-colors border border-[var(--border-subtle)]"
            >
              Contact Us
            </a>
          </div>
        </div>

        {/* What's Included Section */}
        <div className="mb-20">
          <h2 className="text-2xl sm:text-3xl font-bold text-[var(--text-primary)] mb-8 text-center">
            What's Included in All Plans
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 max-w-5xl mx-auto">
            <div className="text-center p-4">
              <div className="text-[var(--text-secondary)]">Git-based deployments</div>
            </div>
            <div className="text-center p-4">
              <div className="text-[var(--text-secondary)]">Zero DevOps setup</div>
            </div>
            <div className="text-center p-4">
              <div className="text-[var(--text-secondary)]">Managed containers</div>
            </div>
            <div className="text-center p-4">
              <div className="text-[var(--text-secondary)]">Automatic restarts</div>
            </div>
            <div className="text-center p-4">
              <div className="text-[var(--text-secondary)]">Secure HTTPS URLs</div>
            </div>
            <div className="text-center p-4">
              <div className="text-[var(--text-secondary)]">Simple dashboard</div>
            </div>
          </div>
        </div>

        {/* What We Don't Charge For Section */}
        <div className="mb-20 p-8 rounded-lg border border-[var(--border-subtle)] bg-[var(--surface)]">
          <h2 className="text-2xl sm:text-3xl font-bold text-[var(--text-primary)] mb-6 text-center">
            What Stackyn Does NOT Charge For
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 max-w-4xl mx-auto">
            <div className="text-center">
              <div className="text-[var(--text-secondary)]">Bandwidth (within fair use)</div>
            </div>
            <div className="text-center">
              <div className="text-[var(--text-secondary)]">SSL certificates</div>
            </div>
            <div className="text-center">
              <div className="text-[var(--text-secondary)]">Redeploys</div>
            </div>
            <div className="text-center">
              <div className="text-[var(--text-secondary)]">App restarts</div>
            </div>
          </div>
          <p className="text-center text-[var(--text-primary)] font-semibold mt-6">No hidden fees.</p>
        </div>

        {/* FAQ Section */}
        <div className="mb-20">
          <h2 className="text-2xl sm:text-3xl font-bold text-[var(--text-primary)] text-center mb-12">
            Frequently Asked Questions
          </h2>
          <div className="max-w-3xl mx-auto space-y-8">
            <div>
              <h3 className="text-xl font-semibold text-[var(--text-primary)] mb-3">What counts as an app?</h3>
              <p className="text-[var(--text-secondary)]">
                One Git repository deployed as a single running service.
              </p>
            </div>
            <div>
              <h3 className="text-xl font-semibold text-[var(--text-primary)] mb-3">Can I upgrade or downgrade anytime?</h3>
              <p className="text-[var(--text-secondary)]">
                Yes. Plans are flexible and can be changed at any time.
              </p>
            </div>
            <div>
              <h3 className="text-xl font-semibold text-[var(--text-primary)] mb-3">What happens if I exceed my limits?</h3>
              <p className="text-[var(--text-secondary)]">
                We'll notify you and help you upgrade — your app won't be shut down suddenly.
              </p>
            </div>
            <div>
              <h3 className="text-xl font-semibold text-[var(--text-primary)] mb-3">Is this suitable for production?</h3>
              <p className="text-[var(--text-secondary)]">
                Yes. Stackyn is built for real-world backend apps with reliability and isolation.
              </p>
            </div>
          </div>
        </div>

        {/* Final CTA Section */}
        <div className="text-center mb-20 p-12 rounded-lg border border-[var(--border-subtle)] bg-[var(--surface)]">
          <h2 className="text-3xl sm:text-4xl font-bold text-[var(--text-primary)] mb-6">
            Start shipping, not configuring servers.
          </h2>
          <a
            href={user ? "https://console.staging.stackyn.com/" : "/login"}
            onClick={handleSignInClick}
            className="inline-block bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-[var(--app-bg)] font-semibold py-4 px-8 rounded-lg transition-colors text-lg mb-4"
          >
            Get Started Free
          </a>
          <p className="text-sm text-[var(--text-muted)]">No credit card required</p>
        </div>

        {/* Founder Note */}
        <div className="max-w-3xl mx-auto p-8 rounded-lg border-l-4 border-[var(--primary)] bg-[var(--primary-muted)]/20">
          <p className="text-[var(--text-primary)] italic text-lg">
            "Stackyn is built for developers who want clarity, speed, and control — without becoming DevOps engineers."
          </p>
        </div>
      </div>

      {/* Footer */}
      <footer className="bg-[var(--sidebar)] text-[var(--text-muted)] py-12 mt-16 border-t border-[var(--border-subtle)]">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4">
            <div className="text-sm">
              <p>&copy; {new Date().getFullYear()} Stackyn. All rights reserved.</p>
            </div>
            <div className="flex items-center gap-6 text-sm">
              <Link to="/terms" className="hover:text-[var(--text-primary)] transition-colors">
                Terms of Service
              </Link>
              <Link to="/privacy" className="hover:text-[var(--text-primary)] transition-colors">
                Privacy Policy
              </Link>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
