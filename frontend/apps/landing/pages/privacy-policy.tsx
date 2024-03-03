import { Title, Text, Container } from '@mantine/core';
import { AppLayout } from '../src/layouts';

export default function PrivacyPolicyPage(): JSX.Element {
  return (
    <AppLayout title="Privacy Policy">
      <Container size="xl">
        <Title order={1}>NOTE: THIS IS A PLACEHOLDER MADE BY ChatGPT! IGNORE THIS!</Title>

        <Text>
          Sure! Here&apos;s an example privacy policy for a cooking app: Privacy Policy for Cooking App Effective Date:
          [Date]
          <p>
            This Privacy Policy describes how our cooking app (&quot;App&quot;) collects, uses, and discloses your
            personal information when you use our services. By using the App, you consent to the practices described in
            this Privacy Policy.
          </p>
          Information We Collect 1.1 Personal Information:
          <p>
            We may collect personal information that you provide directly to us when you use the App, such as your name,
            email address, and profile picture. We only collect personal information that is necessary to provide our
            services and improve the user experience.
          </p>
          1.2 Usage Information:
          <p>
            We may collect information about your usage of the App, including your interactions with recipes,
            preferences, and search queries. This information helps us personalize your experience and improve our
            services.
          </p>
          1.3 Device Information:
          <p>
            We may collect information about the device you use to access the App, including the device type, operating
            system, and unique device identifiers. This information helps us optimize and troubleshoot the App for
            different devices.
          </p>
          Use of Information 2.1 Provide and Improve Services:
          <p>
            We use the information collected to provide and improve our services, including offering personalized
            recipes, recommending relevant content, and enhancing the overall user experience.
          </p>
          2.2 Communication:
          <p>
            We may use your email address to send you important updates, such as changes to our terms or privacy policy,
            or to respond to your inquiries and feedback.
          </p>
          2.3 Analytics:
          <p>
            We may use aggregated and anonymized information for analytical purposes to understand user behavior,
            trends, and preferences. This helps us improve our services and develop new features.
          </p>
          Sharing of Information 3.1 Service Providers:
          <p>
            We may share your information with trusted third-party service providers who assist us in delivering and
            maintaining the App. These providers are obligated to protect your information and use it only for the
            purpose of providing the requested services.
          </p>
          3.2 Legal Requirements:
          <p>
            We may disclose your information if required to do so by law or in response to valid requests by public
            authorities (e.g., court orders or government agencies).
          </p>
          3.3 Business Transfers:
          <p>
            In the event of a merger, acquisition, or sale of our company, your information may be transferred to the
            acquiring entity. We will notify you before your personal information becomes subject to a different privacy
            policy.
          </p>
          Data Security
          <p>
            We take appropriate security measures to protect your personal information against unauthorized access,
            alteration, disclosure, or destruction. However, no method of transmission over the internet or electronic
            storage is completely secure, so we cannot guarantee absolute security.
          </p>
          Children&apos;s Privacy
          <p>
            The App is not intended for use by children under the age of 13. We do not knowingly collect personal
            information from children under 13. If you believe that a child has provided us with their personal
            information, please contact us, and we will promptly delete it.
          </p>
          Changes to this Privacy Policy
          <p>
            We may update this Privacy Policy from time to time. We will notify you of any significant changes by
            posting the new policy on our website or within the App. Your continued use of the App after such
            modifications constitutes your acknowledgment and acceptance of the updated Privacy Policy.
          </p>
          Contact Us
          <p>
            If you have any questions or concerns about this Privacy Policy, please contact us at [contact email
            address].
          </p>
          <p>
            Please note that this is just an example and may not cover all the specific requirements of your cooking
            app. It&apos;s recommended to consult with a legal professional to ensure your privacy policy complies with
            applicable laws and regulations.
          </p>
        </Text>
      </Container>
    </AppLayout>
  );
}
