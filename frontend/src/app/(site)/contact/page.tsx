import { ContactForm } from "@/components/site/ContactForm";

export const dynamic = "force-static";

export const metadata = {
  title: "Contact",
  description: "Get in touch with the Tesdevops team.",
};

export default function ContactPage() {
  return (
    <div className="mx-auto w-full max-w-2xl px-4 py-16 sm:px-6">
      <header className="space-y-3">
        <p className="text-sm uppercase tracking-wide text-muted-foreground">
          Contact
        </p>
        <h1 className="text-4xl font-semibold tracking-tight">
          Tell us what you're building.
        </h1>
        <p className="text-lg leading-relaxed text-muted-foreground">
          The form is validated client-side with Zod and sends a structured log
          to the API on submit. Hook it up to your favourite mailbox or CRM
          when you're ready.
        </p>
      </header>

      <div className="mt-10">
        <ContactForm />
      </div>
    </div>
  );
}
