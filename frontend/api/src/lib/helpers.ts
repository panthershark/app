// BONKERS (Also really cool!): https://www.petermorlion.com/iterating-a-typescript-enum/
export const enumKeys = <O extends object, K extends keyof O = keyof O>(obj: O): K[] => {
	return Object.keys(obj).filter((k) => Number.isNaN(+k)) as K[];
};

export const enumLabel = <O extends object, K extends keyof O = keyof O>(obj: O, val: any): K | undefined => {
	const keys = enumKeys(obj) as K[];
	return keys.find((k) => obj[k] === val);
};

export const enumValue = <O extends object, K extends keyof O = keyof O>(obj: O, key: string) => {
	return obj[key as K];
};
