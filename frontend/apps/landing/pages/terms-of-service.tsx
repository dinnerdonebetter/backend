import { Title, Text, Container } from '@mantine/core';
import { AppLayout } from '../src/layouts';

export default function PrivacyPolicyPage(): JSX.Element {
  return (
    <AppLayout title="Terms of Service">
      <Container size="xl">
        <Title order={1}>NOTE: THIS IS A PLACEHOLDER MADE BY ChatGPT! IGNORE THIS!</Title>

        <Text>
          Effective Date: [Date] Please read these Terms of Service (&quot;Terms&quot;) carefully before using our
          cooking app (&quot;App&quot;). By accessing or using the App, you agree to be bound by these Terms. If you do
          not agree with any part of these Terms, you may not use the App. Use of the App 1.1 Eligibility:
          <p>
            You must be at least 13 years old to use the App. By using the App, you represent and warrant that you are
            13 years of age or older. If you are under the age of 18, you must have the permission of a parent or
            guardian to use the App.
          </p>
          1.2 License:
          <p>
            Subject to your compliance with these Terms, we grant you a limited, non-exclusive, non-transferable, and
            revocable license to use the App for personal, non-commercial purposes. You may not modify, distribute, or
            create derivative works based on the App.
          </p>
          1.3 User Conduct:
          <p>
            You agree to use the App in a responsible and lawful manner. You must not engage in any activity that
            interferes with or disrupts the operation of the App or violates any laws or regulations.
          </p>
          User Content 2.1 Content Ownership:
          <p>
            You retain ownership of any content you submit, upload, or post on the App (&quot;User Content&quot;). By
            submitting User Content, you grant us a non-exclusive, royalty-free, worldwide, perpetual, and irrevocable
            license to use, reproduce, distribute, modify, adapt, publicly display, and perform the User Content in
            connection with the App.
          </p>
          2.2 Content Guidelines:
          <p>
            You are solely responsible for the User Content you submit. You must ensure that your User Content does not
            infringe upon the rights of any third party, including copyrights, trademarks, or privacy rights. You must
            not submit any offensive, illegal, or inappropriate content.
          </p>
          2.3 Monitoring and Removal:
          <p>
            We have the right to monitor User Content and may remove or modify any content that violates these Terms or
            is otherwise objectionable. We are not responsible for any loss, damage, or liability arising from the
            removal or modification of User Content.
          </p>
          Intellectual Property Rights
          <p>
            The App and its content, including but not limited to text, graphics, images, logos, and software, are
            protected by intellectual property laws. You may not use, copy, reproduce, modify, distribute, or create
            derivative works of the App or its content without our prior written consent.
          </p>
          Disclaimer of Warranty
          <p>
            The App is provided on an &quot;as is&quot; and &quot;as available&quot; basis without warranties of any
            kind, whether express or implied. We do not warrant that the App will be error-free, secure, or
            uninterrupted. Your use of the App is at your own risk.
          </p>
          Limitation of Liability
          <p>
            To the maximum extent permitted by law, we shall not be liable for any indirect, incidental, special,
            consequential, or punitive damages, or any loss of profits or revenues, whether incurred directly or
            indirectly, arising from your use of the App.
          </p>
          Indemnification
          <p>
            You agree to indemnify and hold us harmless from any claims, damages, liabilities, and expenses (including
            attorneys&apos; fees) arising out of or related to your use of the App, your violation of these Terms, or
            your violation of any rights of a third party.
          </p>
          Termination
          <p>
            We may terminate or suspend your access to the App at any time, with or without cause, and without prior
            notice. Upon termination, your license to use the App will cease, and any provisions of these Terms that
            should reasonably survive termination will remain in effect.
          </p>
        </Text>
      </Container>
    </AppLayout>
  );
}
