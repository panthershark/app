import { expect, test } from 'vitest';
import { enumKeys } from './helpers';

test('enumKeys', () => {
	enum TestEnum {
		Woah = 'woah',
		Eeek = 3,
		oof
	}
	const got = enumKeys(TestEnum);
	expect(got).toEqual(['Woah', 'Eeek', 'oof']);
});
