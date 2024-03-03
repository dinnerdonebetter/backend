import { ServerTiming } from './index';

describe('basic', () => {
  it('should track basic actions', () => {
    const t = new ServerTiming();

    const testEvent = t.addEvent('test', 'testing');
    setTimeout(() => {
      testEvent.end();
    }, 100);

    setTimeout(() => {
      const actualParts = t.headerValue().split(';');
      expect(actualParts.length).toEqual(3);

      const expected = 'test;desc="testing";';
      const descriptionParts = `${actualParts[0]};${actualParts[1]};`;
      expect(descriptionParts).toEqual(expected);

      const finalParts = actualParts[2].split('=');
      expect(finalParts.length).toEqual(2);
      expect(finalParts[0]).toEqual('dur');
      const actual = parseInt(finalParts[1], 10);

      expect(actual).toBeGreaterThanOrEqual(100);
    }, 1000);
  });
});
