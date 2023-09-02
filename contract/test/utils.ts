export const testFuryAddrs = [
  "fury1esagqd83rhqdtpy5sxhklaxgn58k2m3s3mnpea",
  "fury1mq9qxlhze029lm0frzw2xr6hem8c3k9ts54w0w",
  "fury16g8lzm86f5wwf3x3t67qrpd46sjdpxpfazskwg",
  "fury1wn74shl496ktcfgqsc6yf0vvenhgq0hwuw6z2a",
];

export const bech32Prefix = "fury";

export function tokens(
  value: bigint | number,
  decimals: bigint | number = 18
): bigint {
  return BigInt(value) * 10n ** BigInt(decimals);
}
