import { Component, ComponentProps, JSX } from 'solid-js'

export const PeatchIcon: Component<ComponentProps<'svg'>> = props => {
	return (
		<svg
			{...props}
			viewBox="0 200"
			fill="none"
			xmlns="http://www.w3.org/2000/svg"
		>
			<g filter="url(#filter0_ii_2909_9537)">
				<ellipse cx="100" cy="100" rx="100" ry="100" fill="#FF8C42" />
				<ellipse
					cx="100"
					cy="100"
					rx="100"
					ry="100"
					fill="url(#paint0_radial_2909_9537)"
				/>
				<ellipse
					cx="100"
					cy="100"
					rx="100"
					ry="100"
					fill="url(#paint1_radial_2909_9537)"
				/>
				<ellipse
					cx="100"
					cy="100"
					rx="100"
					ry="100"
					fill="url(#paint2_radial_2909_9537)"
					fill-opacity="0.6"
				/>
			</g>
			<path
				d="M199 100C199 154.676 199 100 199C45.3238 199 1 154.676 1 100C1 45.3238 1 100 1C154.676 1 199 45.3238 199 100Z"
				stroke="url(#paint3_linear_2909_9537)"
				stroke-opacity="0.16"
				stroke-width="2"
			/>
			<defs>
				<filter
					id="filter0_ii_2909_9537"
					x="-10.5"
					y="0"
					width="210.5"
					height="223"
					filterUnits="userSpaceOnUse"
					color-interpolation-filters="sRGB"
				>
					<feFlood flood-opacity="0" result="BackgroundImageFix" />
					<feBlend
						mode="normal"
						in="SourceGraphic"
						in2="BackgroundImageFix"
						result="shape"
					/>
					<feColorMatrix
						in="SourceAlpha"
						type="matrix"
						values="0 127 0"
						result="hardAlpha"
					/>
					<feOffset />
					<feGaussianBlur stdDeviation="2.5" />
					<feComposite in2="hardAlpha" operator="arithmetic" k2="-1" k3="1" />
					<feColorMatrix
						type="matrix"
						values="0 1 0 1 0 1 0 0.45 0"
					/>
					<feBlend
						mode="normal"
						in2="shape"
						result="effect1_innerShadow_2909_9537"
					/>
					<feColorMatrix
						in="SourceAlpha"
						type="matrix"
						values="0 127 0"
						result="hardAlpha"
					/>
					<feMorphology
						radius="2.5"
						operator="dilate"
						in="SourceAlpha"
						result="effect2_innerShadow_2909_9537"
					/>
					<feOffset dx="-10.5" dy="23" />
					<feGaussianBlur stdDeviation="13.75" />
					<feComposite in2="hardAlpha" operator="arithmetic" k2="-1" k3="1" />
					<feColorMatrix
						type="matrix"
						values="0 1 0 1 0 1 0 0.55 0"
					/>
					<feBlend
						mode="normal"
						in2="effect1_innerShadow_2909_9537"
						result="effect2_innerShadow_2909_9537"
					/>
				</filter>
				<radialGradient
					id="paint0_radial_2909_9537"
					cx="0"
					cy="0"
					r="1"
					gradientUnits="userSpaceOnUse"
					gradientTransform="translate(168.571 162.857) rotate(-159.775) scale(115.705 115.705)"
				>
					<stop stop-color="#F35D28" />
					<stop offset="1" stop-color="#F35D28" stop-opacity="0" />
				</radialGradient>
				<radialGradient
					id="paint1_radial_2909_9537"
					cx="0"
					cy="0"
					r="1"
					gradientUnits="userSpaceOnUse"
					gradientTransform="translate(52.8571 200) rotate(-59.0362) scale(141.609 141.609)"
				>
					<stop stop-color="#FFD67E" />
					<stop offset="1" stop-color="#FFD77F" stop-opacity="0" />
				</radialGradient>
				<radialGradient
					id="paint2_radial_2909_9537"
					cx="0"
					cy="0"
					r="1"
					gradientUnits="userSpaceOnUse"
					gradientTransform="translate(130 71.4286) rotate(100.437) scale(165.597 132.187)"
				>
					<stop stop-color="white" />
					<stop offset="0.489583" stop-color="white" stop-opacity="0" />
				</radialGradient>
				<linearGradient
					id="paint3_linear_2909_9537"
					x1="100"
					y1="0"
					x2="100"
					y2="200"
					gradientUnits="userSpaceOnUse"
				>
					<stop stop-color="#FFB594" />
					<stop offset="1" stop-color="white" />
				</linearGradient>
			</defs>
		</svg>
	)
}
