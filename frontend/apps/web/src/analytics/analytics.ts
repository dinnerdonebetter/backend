import { Analytics } from '@segment/analytics-node';
import { AnalyticsBrowser } from '@segment/analytics-next';
import { GetServerSidePropsContext } from 'next/types';

import { AnalyticsEvent, PageName } from './events';

class serverAnalyticsWrapper {
  noopMode: boolean = false;
  analytics?: Analytics;

  constructor() {
    const apiKey = process.env.NEXT_PUBLIC_SEGMENT_API_TOKEN || '';
    if (apiKey === '') {
      this.noopMode = true;
    } else {
      this.analytics = new Analytics({ writeKey: apiKey });
    }
  }

  track(userID: string, event: AnalyticsEvent, properties: Record<string, any>) {
    if (!this.noopMode && userID.trim() !== '') {
      this.analytics?.track({ event, properties, userId: userID });
    }
  }

  page(userID: string, pageName: PageName, context: GetServerSidePropsContext, properties: Record<string, any>) {
    if (!this.noopMode && userID.trim() !== '') {
      if (context.query) {
        properties.query = context.query;
      }

      if (context.resolvedUrl) {
        properties.path = context.resolvedUrl;
      }

      if (context.locale) {
        properties.locale = context.locale;
      }

      if (context.params) {
        properties.params = context.params;
      }

      this.analytics?.page({ name: pageName, properties, userId: userID });
    }
  }

  identify(userID: string = '', traits: Record<string, any>) {
    if (!this.noopMode && userID.trim() !== '') {
      this.analytics?.identify({ userId: userID, traits });
    }
  }

  group(userID: string, groupId: string, traits: Record<string, any> = {}) {
    if (!this.noopMode && userID.trim() !== '') {
      this.analytics?.group({ userId: userID, groupId, traits });
    }
  }
}

class browserAnalyticsWrapper {
  noopMode: boolean = false;
  analytics?: AnalyticsBrowser;

  constructor() {
    const apiKey = process.env.NEXT_PUBLIC_SEGMENT_API_TOKEN || '';
    if (apiKey === '') {
      this.noopMode = true;
    } else {
      this.analytics = AnalyticsBrowser.load({ writeKey: apiKey });
    }
  }

  track(event: AnalyticsEvent, properties: Record<string, any>) {
    if (!this.noopMode) {
      this.analytics?.track(event, properties);
    }
  }

  page(name: string, properties: Record<string, any>) {
    if (!this.noopMode) {
      this.analytics?.page({ name, properties });
    }
  }

  identify(userID: string, traits: Record<string, any>) {
    if (!this.noopMode) {
      this.analytics?.identify(userID, { traits });
    }
  }
}

export const browserSideAnalytics = new browserAnalyticsWrapper();
export const serverSideAnalytics = new serverAnalyticsWrapper();
