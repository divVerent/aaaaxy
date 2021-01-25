<?xml version="1.0" encoding="UTF-8"?>
<tileset version="1.4" tiledversion="1.4.3" name="tiles" tilewidth="16" tileheight="16" tilecount="192" columns="16" objectalignment="topleft">
 <grid orientation="orthogonal" width="1" height="1"/>
 <terraintypes>
  <terrain name="Block" tile="9"/>
  <terrain name="Open" tile="63"/>
 </terraintypes>
 <tile id="0">
  <properties>
   <property name="img.EN" value="notile.png"/>
   <property name="img.ES" value="notile.png"/>
   <property name="img.NE" value="notile.png"/>
   <property name="img.NW" value="notile.png"/>
   <property name="img.SE" value="notile.png"/>
   <property name="img.SW" value="notile.png"/>
   <property name="img.WN" value="notile.png"/>
   <property name="img.WS" value="notile.png"/>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="notile.png"/>
 </tile>
 <tile id="1" type="WarpZone">
  <properties>
   <property name="orientation" value="EN"/>
   <property name="type" value="WarpZone"/>
  </properties>
  <image width="16" height="16" source="warpzone_en.png"/>
 </tile>
 <tile id="2" type="WarpZone">
  <properties>
   <property name="orientation" value="ES"/>
   <property name="type" value="WarpZone"/>
  </properties>
  <image width="16" height="16" source="warpzone_es.png"/>
 </tile>
 <tile id="3" type="WarpZone">
  <properties>
   <property name="orientation" value="NE"/>
   <property name="type" value="WarpZone"/>
  </properties>
  <image width="16" height="16" source="warpzone_ne.png"/>
 </tile>
 <tile id="4" type="WarpZone">
  <properties>
   <property name="orientation" value="NW"/>
   <property name="type" value="WarpZone"/>
  </properties>
  <image width="16" height="16" source="warpzone_nw.png"/>
 </tile>
 <tile id="5" type="WarpZone">
  <properties>
   <property name="orientation" value="SE"/>
   <property name="type" value="WarpZone"/>
  </properties>
  <image width="16" height="16" source="warpzone_se.png"/>
 </tile>
 <tile id="6" type="WarpZone">
  <properties>
   <property name="orientation" value="SW"/>
   <property name="type" value="WarpZone"/>
  </properties>
  <image width="16" height="16" source="warpzone_sw.png"/>
 </tile>
 <tile id="7" type="WarpZone">
  <properties>
   <property name="orientation" value="WN"/>
   <property name="type" value="WarpZone"/>
  </properties>
  <image width="16" height="16" source="warpzone_wn.png"/>
 </tile>
 <tile id="8" type="WarpZone">
  <properties>
   <property name="orientation" value="WS"/>
   <property name="type" value="WarpZone"/>
  </properties>
  <image width="16" height="16" source="warpzone_ws.png"/>
 </tile>
 <tile id="9">
  <properties>
   <property name="img.EN" value="wall_.png"/>
   <property name="img.ES" value="wall_.png"/>
   <property name="img.NE" value="wall_.png"/>
   <property name="img.NW" value="wall_.png"/>
   <property name="img.SE" value="wall_.png"/>
   <property name="img.SW" value="wall_.png"/>
   <property name="img.WN" value="wall_.png"/>
   <property name="img.WS" value="wall_.png"/>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_.png"/>
 </tile>
 <tile id="10">
  <properties>
   <property name="img.EN" value="wall_w.png"/>
   <property name="img.ES" value="wall_w.png"/>
   <property name="img.NE" value="wall_s.png"/>
   <property name="img.NW" value="wall_s.png"/>
   <property name="img.SE" value="wall_n.png"/>
   <property name="img.SW" value="wall_n.png"/>
   <property name="img.WN" value="wall_e.png"/>
   <property name="img.WS" value="wall_e.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_w.png"/>
 </tile>
 <tile id="11">
  <properties>
   <property name="img.EN" value="wall_n.png"/>
   <property name="img.ES" value="wall_s.png"/>
   <property name="img.NE" value="wall_e.png"/>
   <property name="img.NW" value="wall_w.png"/>
   <property name="img.SE" value="wall_e.png"/>
   <property name="img.SW" value="wall_w.png"/>
   <property name="img.WN" value="wall_n.png"/>
   <property name="img.WS" value="wall_s.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_s.png"/>
 </tile>
 <tile id="12">
  <properties>
   <property name="img.EN" value="wall_n_w_nw0.png"/>
   <property name="img.ES" value="wall_s_w_sw0.png"/>
   <property name="img.NE" value="wall_e_s_se0.png"/>
   <property name="img.NW" value="wall_s_w_sw0.png"/>
   <property name="img.SE" value="wall_n_e_ne0.png"/>
   <property name="img.SW" value="wall_n_w_nw0.png"/>
   <property name="img.WN" value="wall_n_e_ne0.png"/>
   <property name="img.WS" value="wall_e_s_se0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_s_w_sw0.png"/>
 </tile>
 <tile id="13" terrain="1,1,0,1">
  <properties>
   <property name="img.EN" value="wall_n_w_nw1.png"/>
   <property name="img.ES" value="wall_s_w_sw1.png"/>
   <property name="img.NE" value="wall_e_s_se1.png"/>
   <property name="img.NW" value="wall_s_w_sw1.png"/>
   <property name="img.SE" value="wall_n_e_ne1.png"/>
   <property name="img.SW" value="wall_n_w_nw1.png"/>
   <property name="img.WN" value="wall_n_e_ne1.png"/>
   <property name="img.WS" value="wall_e_s_se1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_s_w_sw1.png"/>
 </tile>
 <tile id="14">
  <properties>
   <property name="img.EN" value="wall_e.png"/>
   <property name="img.ES" value="wall_e.png"/>
   <property name="img.NE" value="wall_n.png"/>
   <property name="img.NW" value="wall_n.png"/>
   <property name="img.SE" value="wall_s.png"/>
   <property name="img.SW" value="wall_s.png"/>
   <property name="img.WN" value="wall_w.png"/>
   <property name="img.WS" value="wall_w.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_e.png"/>
 </tile>
 <tile id="15">
  <properties>
   <property name="img.EN" value="wall_e_w.png"/>
   <property name="img.ES" value="wall_e_w.png"/>
   <property name="img.NE" value="wall_n_s.png"/>
   <property name="img.NW" value="wall_n_s.png"/>
   <property name="img.SE" value="wall_n_s.png"/>
   <property name="img.SW" value="wall_n_s.png"/>
   <property name="img.WN" value="wall_e_w.png"/>
   <property name="img.WS" value="wall_e_w.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_e_w.png"/>
 </tile>
 <tile id="16">
  <properties>
   <property name="img.EN" value="wall_n_e_ne0.png"/>
   <property name="img.ES" value="wall_e_s_se0.png"/>
   <property name="img.NE" value="wall_n_e_ne0.png"/>
   <property name="img.NW" value="wall_n_w_nw0.png"/>
   <property name="img.SE" value="wall_e_s_se0.png"/>
   <property name="img.SW" value="wall_s_w_sw0.png"/>
   <property name="img.WN" value="wall_n_w_nw0.png"/>
   <property name="img.WS" value="wall_s_w_sw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_e_s_se0.png"/>
 </tile>
 <tile id="17" terrain="1,1,1,0">
  <properties>
   <property name="img.EN" value="wall_n_e_ne1.png"/>
   <property name="img.ES" value="wall_e_s_se1.png"/>
   <property name="img.NE" value="wall_n_e_ne1.png"/>
   <property name="img.NW" value="wall_n_w_nw1.png"/>
   <property name="img.SE" value="wall_e_s_se1.png"/>
   <property name="img.SW" value="wall_s_w_sw1.png"/>
   <property name="img.WN" value="wall_n_w_nw1.png"/>
   <property name="img.WS" value="wall_s_w_sw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_e_s_se1.png"/>
 </tile>
 <tile id="18">
  <properties>
   <property name="img.EN" value="wall_n_e_w_ne0_nw0.png"/>
   <property name="img.ES" value="wall_e_s_w_se0_sw0.png"/>
   <property name="img.NE" value="wall_n_e_s_ne0_se0.png"/>
   <property name="img.NW" value="wall_n_s_w_sw0_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_ne0_se0.png"/>
   <property name="img.SW" value="wall_n_s_w_sw0_nw0.png"/>
   <property name="img.WN" value="wall_n_e_w_ne0_nw0.png"/>
   <property name="img.WS" value="wall_e_s_w_se0_sw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_e_s_w_se0_sw0.png"/>
 </tile>
 <tile id="19">
  <properties>
   <property name="img.EN" value="wall_n_e_w_ne0_nw1.png"/>
   <property name="img.ES" value="wall_e_s_w_se0_sw1.png"/>
   <property name="img.NE" value="wall_n_e_s_ne0_se1.png"/>
   <property name="img.NW" value="wall_n_s_w_sw1_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_ne1_se0.png"/>
   <property name="img.SW" value="wall_n_s_w_sw0_nw1.png"/>
   <property name="img.WN" value="wall_n_e_w_ne1_nw0.png"/>
   <property name="img.WS" value="wall_e_s_w_se1_sw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_e_s_w_se0_sw1.png"/>
 </tile>
 <tile id="20">
  <properties>
   <property name="img.EN" value="wall_n_e_w_ne1_nw0.png"/>
   <property name="img.ES" value="wall_e_s_w_se1_sw0.png"/>
   <property name="img.NE" value="wall_n_e_s_ne1_se0.png"/>
   <property name="img.NW" value="wall_n_s_w_sw0_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_ne0_se1.png"/>
   <property name="img.SW" value="wall_n_s_w_sw1_nw0.png"/>
   <property name="img.WN" value="wall_n_e_w_ne0_nw1.png"/>
   <property name="img.WS" value="wall_e_s_w_se0_sw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_e_s_w_se1_sw0.png"/>
 </tile>
 <tile id="21" terrain="1,1,0,0">
  <properties>
   <property name="img.EN" value="wall_n_e_w_ne1_nw1.png"/>
   <property name="img.ES" value="wall_e_s_w_se1_sw1.png"/>
   <property name="img.NE" value="wall_n_e_s_ne1_se1.png"/>
   <property name="img.NW" value="wall_n_s_w_sw1_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_ne1_se1.png"/>
   <property name="img.SW" value="wall_n_s_w_sw1_nw1.png"/>
   <property name="img.WN" value="wall_n_e_w_ne1_nw1.png"/>
   <property name="img.WS" value="wall_e_s_w_se1_sw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_e_s_w_se1_sw1.png"/>
 </tile>
 <tile id="22">
  <properties>
   <property name="img.EN" value="wall_s.png"/>
   <property name="img.ES" value="wall_n.png"/>
   <property name="img.NE" value="wall_w.png"/>
   <property name="img.NW" value="wall_e.png"/>
   <property name="img.SE" value="wall_w.png"/>
   <property name="img.SW" value="wall_e.png"/>
   <property name="img.WN" value="wall_s.png"/>
   <property name="img.WS" value="wall_n.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n.png"/>
 </tile>
 <tile id="23">
  <properties>
   <property name="img.EN" value="wall_s_w_sw0.png"/>
   <property name="img.ES" value="wall_n_w_nw0.png"/>
   <property name="img.NE" value="wall_s_w_sw0.png"/>
   <property name="img.NW" value="wall_e_s_se0.png"/>
   <property name="img.SE" value="wall_n_w_nw0.png"/>
   <property name="img.SW" value="wall_n_e_ne0.png"/>
   <property name="img.WN" value="wall_e_s_se0.png"/>
   <property name="img.WS" value="wall_n_e_ne0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_w_nw0.png"/>
 </tile>
 <tile id="24" terrain="0,1,1,1">
  <properties>
   <property name="img.EN" value="wall_s_w_sw1.png"/>
   <property name="img.ES" value="wall_n_w_nw1.png"/>
   <property name="img.NE" value="wall_s_w_sw1.png"/>
   <property name="img.NW" value="wall_e_s_se1.png"/>
   <property name="img.SE" value="wall_n_w_nw1.png"/>
   <property name="img.SW" value="wall_n_e_ne1.png"/>
   <property name="img.WN" value="wall_e_s_se1.png"/>
   <property name="img.WS" value="wall_n_e_ne1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_w_nw1.png"/>
 </tile>
 <tile id="25">
  <properties>
   <property name="img.EN" value="wall_n_s.png"/>
   <property name="img.ES" value="wall_n_s.png"/>
   <property name="img.NE" value="wall_e_w.png"/>
   <property name="img.NW" value="wall_e_w.png"/>
   <property name="img.SE" value="wall_e_w.png"/>
   <property name="img.SW" value="wall_e_w.png"/>
   <property name="img.WN" value="wall_n_s.png"/>
   <property name="img.WS" value="wall_n_s.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_s.png"/>
 </tile>
 <tile id="26">
  <properties>
   <property name="img.EN" value="wall_n_s_w_sw0_nw0.png"/>
   <property name="img.ES" value="wall_n_s_w_sw0_nw0.png"/>
   <property name="img.NE" value="wall_e_s_w_se0_sw0.png"/>
   <property name="img.NW" value="wall_e_s_w_se0_sw0.png"/>
   <property name="img.SE" value="wall_n_e_w_ne0_nw0.png"/>
   <property name="img.SW" value="wall_n_e_w_ne0_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_ne0_se0.png"/>
   <property name="img.WS" value="wall_n_e_s_ne0_se0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_s_w_sw0_nw0.png"/>
 </tile>
 <tile id="27">
  <properties>
   <property name="img.EN" value="wall_n_s_w_sw1_nw0.png"/>
   <property name="img.ES" value="wall_n_s_w_sw0_nw1.png"/>
   <property name="img.NE" value="wall_e_s_w_se0_sw1.png"/>
   <property name="img.NW" value="wall_e_s_w_se1_sw0.png"/>
   <property name="img.SE" value="wall_n_e_w_ne0_nw1.png"/>
   <property name="img.SW" value="wall_n_e_w_ne1_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_ne0_se1.png"/>
   <property name="img.WS" value="wall_n_e_s_ne1_se0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_s_w_sw0_nw1.png"/>
 </tile>
 <tile id="28">
  <properties>
   <property name="img.EN" value="wall_n_s_w_sw0_nw1.png"/>
   <property name="img.ES" value="wall_n_s_w_sw1_nw0.png"/>
   <property name="img.NE" value="wall_e_s_w_se1_sw0.png"/>
   <property name="img.NW" value="wall_e_s_w_se0_sw1.png"/>
   <property name="img.SE" value="wall_n_e_w_ne1_nw0.png"/>
   <property name="img.SW" value="wall_n_e_w_ne0_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_ne1_se0.png"/>
   <property name="img.WS" value="wall_n_e_s_ne0_se1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_s_w_sw1_nw0.png"/>
 </tile>
 <tile id="29" terrain="0,1,0,1">
  <properties>
   <property name="img.EN" value="wall_n_s_w_sw1_nw1.png"/>
   <property name="img.ES" value="wall_n_s_w_sw1_nw1.png"/>
   <property name="img.NE" value="wall_e_s_w_se1_sw1.png"/>
   <property name="img.NW" value="wall_e_s_w_se1_sw1.png"/>
   <property name="img.SE" value="wall_n_e_w_ne1_nw1.png"/>
   <property name="img.SW" value="wall_n_e_w_ne1_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_ne1_se1.png"/>
   <property name="img.WS" value="wall_n_e_s_ne1_se1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_s_w_sw1_nw1.png"/>
 </tile>
 <tile id="30">
  <properties>
   <property name="img.EN" value="wall_e_s_se0.png"/>
   <property name="img.ES" value="wall_n_e_ne0.png"/>
   <property name="img.NE" value="wall_n_w_nw0.png"/>
   <property name="img.NW" value="wall_n_e_ne0.png"/>
   <property name="img.SE" value="wall_s_w_sw0.png"/>
   <property name="img.SW" value="wall_e_s_se0.png"/>
   <property name="img.WN" value="wall_s_w_sw0.png"/>
   <property name="img.WS" value="wall_n_w_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_ne0.png"/>
 </tile>
 <tile id="31" terrain="1,0,1,1">
  <properties>
   <property name="img.EN" value="wall_e_s_se1.png"/>
   <property name="img.ES" value="wall_n_e_ne1.png"/>
   <property name="img.NE" value="wall_n_w_nw1.png"/>
   <property name="img.NW" value="wall_n_e_ne1.png"/>
   <property name="img.SE" value="wall_s_w_sw1.png"/>
   <property name="img.SW" value="wall_e_s_se1.png"/>
   <property name="img.WN" value="wall_s_w_sw1.png"/>
   <property name="img.WS" value="wall_n_w_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_ne1.png"/>
 </tile>
 <tile id="32">
  <properties>
   <property name="img.EN" value="wall_e_s_w_se0_sw0.png"/>
   <property name="img.ES" value="wall_n_e_w_ne0_nw0.png"/>
   <property name="img.NE" value="wall_n_s_w_sw0_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_ne0_se0.png"/>
   <property name="img.SE" value="wall_n_s_w_sw0_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_ne0_se0.png"/>
   <property name="img.WN" value="wall_e_s_w_se0_sw0.png"/>
   <property name="img.WS" value="wall_n_e_w_ne0_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_w_ne0_nw0.png"/>
 </tile>
 <tile id="33">
  <properties>
   <property name="img.EN" value="wall_e_s_w_se0_sw1.png"/>
   <property name="img.ES" value="wall_n_e_w_ne0_nw1.png"/>
   <property name="img.NE" value="wall_n_s_w_sw1_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_ne0_se1.png"/>
   <property name="img.SE" value="wall_n_s_w_sw0_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_ne1_se0.png"/>
   <property name="img.WN" value="wall_e_s_w_se1_sw0.png"/>
   <property name="img.WS" value="wall_n_e_w_ne1_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_w_ne0_nw1.png"/>
 </tile>
 <tile id="34">
  <properties>
   <property name="img.EN" value="wall_e_s_w_se1_sw0.png"/>
   <property name="img.ES" value="wall_n_e_w_ne1_nw0.png"/>
   <property name="img.NE" value="wall_n_s_w_sw0_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_ne1_se0.png"/>
   <property name="img.SE" value="wall_n_s_w_sw1_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_ne0_se1.png"/>
   <property name="img.WN" value="wall_e_s_w_se0_sw1.png"/>
   <property name="img.WS" value="wall_n_e_w_ne0_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_w_ne1_nw0.png"/>
 </tile>
 <tile id="35" terrain="0,0,1,1">
  <properties>
   <property name="img.EN" value="wall_e_s_w_se1_sw1.png"/>
   <property name="img.ES" value="wall_n_e_w_ne1_nw1.png"/>
   <property name="img.NE" value="wall_n_s_w_sw1_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_ne1_se1.png"/>
   <property name="img.SE" value="wall_n_s_w_sw1_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_ne1_se1.png"/>
   <property name="img.WN" value="wall_e_s_w_se1_sw1.png"/>
   <property name="img.WS" value="wall_n_e_w_ne1_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_w_ne1_nw1.png"/>
 </tile>
 <tile id="36">
  <properties>
   <property name="img.EN" value="wall_n_e_s_ne0_se0.png"/>
   <property name="img.ES" value="wall_n_e_s_ne0_se0.png"/>
   <property name="img.NE" value="wall_n_e_w_ne0_nw0.png"/>
   <property name="img.NW" value="wall_n_e_w_ne0_nw0.png"/>
   <property name="img.SE" value="wall_e_s_w_se0_sw0.png"/>
   <property name="img.SW" value="wall_e_s_w_se0_sw0.png"/>
   <property name="img.WN" value="wall_n_s_w_sw0_nw0.png"/>
   <property name="img.WS" value="wall_n_s_w_sw0_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_ne0_se0.png"/>
 </tile>
 <tile id="37">
  <properties>
   <property name="img.EN" value="wall_n_e_s_ne1_se0.png"/>
   <property name="img.ES" value="wall_n_e_s_ne0_se1.png"/>
   <property name="img.NE" value="wall_n_e_w_ne1_nw0.png"/>
   <property name="img.NW" value="wall_n_e_w_ne0_nw1.png"/>
   <property name="img.SE" value="wall_e_s_w_se1_sw0.png"/>
   <property name="img.SW" value="wall_e_s_w_se0_sw1.png"/>
   <property name="img.WN" value="wall_n_s_w_sw0_nw1.png"/>
   <property name="img.WS" value="wall_n_s_w_sw1_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_ne0_se1.png"/>
 </tile>
 <tile id="38">
  <properties>
   <property name="img.EN" value="wall_n_e_s_ne0_se1.png"/>
   <property name="img.ES" value="wall_n_e_s_ne1_se0.png"/>
   <property name="img.NE" value="wall_n_e_w_ne0_nw1.png"/>
   <property name="img.NW" value="wall_n_e_w_ne1_nw0.png"/>
   <property name="img.SE" value="wall_e_s_w_se0_sw1.png"/>
   <property name="img.SW" value="wall_e_s_w_se1_sw0.png"/>
   <property name="img.WN" value="wall_n_s_w_sw1_nw0.png"/>
   <property name="img.WS" value="wall_n_s_w_sw0_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_ne1_se0.png"/>
 </tile>
 <tile id="39" terrain="1,0,1,0">
  <properties>
   <property name="img.EN" value="wall_n_e_s_ne1_se1.png"/>
   <property name="img.ES" value="wall_n_e_s_ne1_se1.png"/>
   <property name="img.NE" value="wall_n_e_w_ne1_nw1.png"/>
   <property name="img.NW" value="wall_n_e_w_ne1_nw1.png"/>
   <property name="img.SE" value="wall_e_s_w_se1_sw1.png"/>
   <property name="img.SW" value="wall_e_s_w_se1_sw1.png"/>
   <property name="img.WN" value="wall_n_s_w_sw1_nw1.png"/>
   <property name="img.WS" value="wall_n_s_w_sw1_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_ne1_se1.png"/>
 </tile>
 <tile id="40">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne0_se0_sw0_nw0.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne0_se0_sw0_nw0.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne0_se0_sw0_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne0_se0_sw0_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne0_se0_sw0_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne0_se0_sw0_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne0_se0_sw0_nw0.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne0_se0_sw0_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne0_se0_sw0_nw0.png"/>
 </tile>
 <tile id="41">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne0_se0_sw1_nw0.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne0_se0_sw0_nw1.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne0_se0_sw1_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne0_se1_sw0_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne0_se0_sw0_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne1_se0_sw0_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne0_se1_sw0_nw0.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne1_se0_sw0_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne0_se0_sw0_nw1.png"/>
 </tile>
 <tile id="42">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne0_se0_sw0_nw1.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne0_se0_sw1_nw0.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne0_se1_sw0_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne0_se0_sw1_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne1_se0_sw0_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne0_se0_sw0_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne1_se0_sw0_nw0.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne0_se1_sw0_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne0_se0_sw1_nw0.png"/>
 </tile>
 <tile id="43">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne0_se0_sw1_nw1.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne0_se0_sw1_nw1.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne0_se1_sw1_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne0_se1_sw1_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne1_se0_sw0_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne1_se0_sw0_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne1_se1_sw0_nw0.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne1_se1_sw0_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne0_se0_sw1_nw1.png"/>
 </tile>
 <tile id="44">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne1_se0_sw0_nw0.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne0_se1_sw0_nw0.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne1_se0_sw0_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne0_se0_sw0_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne0_se1_sw0_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne0_se0_sw1_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne0_se0_sw0_nw1.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne0_se0_sw1_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne0_se1_sw0_nw0.png"/>
 </tile>
 <tile id="45" terrain="0,1,1,0">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne1_se0_sw1_nw0.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne0_se1_sw0_nw1.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne1_se0_sw1_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne0_se1_sw0_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne0_se1_sw0_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne1_se0_sw1_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne0_se1_sw0_nw1.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne1_se0_sw1_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne0_se1_sw0_nw1.png"/>
 </tile>
 <tile id="46">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne1_se0_sw0_nw1.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne0_se1_sw1_nw0.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne1_se1_sw0_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne0_se0_sw1_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne1_se1_sw0_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne0_se0_sw1_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne1_se0_sw0_nw1.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne0_se1_sw1_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne0_se1_sw1_nw0.png"/>
 </tile>
 <tile id="47" terrain="0,1,0,0">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne1_se0_sw1_nw1.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne0_se1_sw1_nw1.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne1_se1_sw1_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne0_se1_sw1_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne1_se1_sw0_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne1_se0_sw1_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne1_se1_sw0_nw1.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne1_se1_sw1_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne0_se1_sw1_nw1.png"/>
 </tile>
 <tile id="48">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne0_se1_sw0_nw0.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne1_se0_sw0_nw0.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne0_se0_sw0_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne1_se0_sw0_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne0_se0_sw1_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne0_se1_sw0_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne0_se0_sw1_nw0.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne0_se0_sw0_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne1_se0_sw0_nw0.png"/>
 </tile>
 <tile id="49">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne0_se1_sw1_nw0.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne1_se0_sw0_nw1.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne0_se0_sw1_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne1_se1_sw0_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne0_se0_sw1_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne1_se1_sw0_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne0_se1_sw1_nw0.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne1_se0_sw0_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne1_se0_sw0_nw1.png"/>
 </tile>
 <tile id="50" terrain="1,0,0,1">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne0_se1_sw0_nw1.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne1_se0_sw1_nw0.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne0_se1_sw0_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne1_se0_sw1_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne1_se0_sw1_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne0_se1_sw0_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne1_se0_sw1_nw0.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne0_se1_sw0_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne1_se0_sw1_nw0.png"/>
 </tile>
 <tile id="51" terrain="0,0,0,1">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne0_se1_sw1_nw1.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne1_se0_sw1_nw1.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne0_se1_sw1_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne1_se1_sw1_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne1_se0_sw1_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne1_se1_sw0_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne1_se1_sw1_nw0.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne1_se1_sw0_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne1_se0_sw1_nw1.png"/>
 </tile>
 <tile id="52">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne1_se1_sw0_nw0.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne1_se1_sw0_nw0.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne1_se0_sw0_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne1_se0_sw0_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne0_se1_sw1_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne0_se1_sw1_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne0_se0_sw1_nw1.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne0_se0_sw1_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne1_se1_sw0_nw0.png"/>
 </tile>
 <tile id="53" terrain="0,0,1,0">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne1_se1_sw1_nw0.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne1_se1_sw0_nw1.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne1_se0_sw1_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne1_se1_sw0_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne0_se1_sw1_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne1_se1_sw1_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne0_se1_sw1_nw1.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne1_se0_sw1_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne1_se1_sw0_nw1.png"/>
 </tile>
 <tile id="54" terrain="1,0,0,0">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne1_se1_sw0_nw1.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne1_se1_sw1_nw0.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne1_se1_sw0_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne1_se0_sw1_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne1_se1_sw1_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne0_se1_sw1_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne1_se0_sw1_nw1.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne0_se1_sw1_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne1_se1_sw1_nw0.png"/>
 </tile>
 <tile id="55" terrain="0,0,0,0">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne1_se1_sw1_nw1.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne1_se1_sw1_nw1.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne1_se1_sw1_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne1_se1_sw1_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne1_se1_sw1_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne1_se1_sw1_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne1_se1_sw1_nw1.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne1_se1_sw1_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne1_se1_sw1_nw1.png"/>
 </tile>
 <tile id="56" terrain="1,1,1,1" probability="0">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_8.png"/>
 </tile>
 <tile id="57" terrain="1,1,1,1" probability="0">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_9.png"/>
 </tile>
 <tile id="58" terrain="1,1,1,1" probability="0">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_a.png"/>
 </tile>
 <tile id="59" terrain="1,1,1,1" probability="0">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_b.png"/>
 </tile>
 <tile id="60" terrain="1,1,1,1" probability="0">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_c.png"/>
 </tile>
 <tile id="61" terrain="1,1,1,1" probability="0">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_d.png"/>
 </tile>
 <tile id="62" terrain="1,1,1,1" probability="0">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_e.png"/>
 </tile>
 <tile id="63" terrain="1,1,1,1">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_f.png"/>
 </tile>
 <tile id="64">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_88_h.png"/>
 </tile>
 <tile id="65">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_88_v.png"/>
 </tile>
 <tile id="66">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_89_h.png"/>
 </tile>
 <tile id="67">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_89_v.png"/>
 </tile>
 <tile id="68">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_8a_h.png"/>
 </tile>
 <tile id="69">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_8a_v.png"/>
 </tile>
 <tile id="70">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_8b_h.png"/>
 </tile>
 <tile id="71">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_8b_v.png"/>
 </tile>
 <tile id="72">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_8c_h.png"/>
 </tile>
 <tile id="73">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_8c_v.png"/>
 </tile>
 <tile id="74">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_8d_h.png"/>
 </tile>
 <tile id="75">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_8d_v.png"/>
 </tile>
 <tile id="76">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_8e_h.png"/>
 </tile>
 <tile id="77">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_8e_v.png"/>
 </tile>
 <tile id="78">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_8f_h.png"/>
 </tile>
 <tile id="79">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_8f_v.png"/>
 </tile>
 <tile id="80">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_98_h.png"/>
 </tile>
 <tile id="81">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_98_v.png"/>
 </tile>
 <tile id="82">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_99_h.png"/>
 </tile>
 <tile id="83">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_99_v.png"/>
 </tile>
 <tile id="84">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_9a_h.png"/>
 </tile>
 <tile id="85">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_9a_v.png"/>
 </tile>
 <tile id="86">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_9b_h.png"/>
 </tile>
 <tile id="87">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_9b_v.png"/>
 </tile>
 <tile id="88">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_9c_h.png"/>
 </tile>
 <tile id="89">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_9c_v.png"/>
 </tile>
 <tile id="90">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_9d_h.png"/>
 </tile>
 <tile id="91">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_9d_v.png"/>
 </tile>
 <tile id="92">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_9e_h.png"/>
 </tile>
 <tile id="93">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_9e_v.png"/>
 </tile>
 <tile id="94">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_9f_h.png"/>
 </tile>
 <tile id="95">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_9f_v.png"/>
 </tile>
 <tile id="96">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_a8_h.png"/>
 </tile>
 <tile id="97">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_a8_v.png"/>
 </tile>
 <tile id="98">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_a9_h.png"/>
 </tile>
 <tile id="99">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_a9_v.png"/>
 </tile>
 <tile id="100">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_aa_h.png"/>
 </tile>
 <tile id="101">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_aa_v.png"/>
 </tile>
 <tile id="102">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ab_h.png"/>
 </tile>
 <tile id="103">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ab_v.png"/>
 </tile>
 <tile id="104">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ac_h.png"/>
 </tile>
 <tile id="105">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ac_v.png"/>
 </tile>
 <tile id="106">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ad_h.png"/>
 </tile>
 <tile id="107">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ad_v.png"/>
 </tile>
 <tile id="108">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ae_h.png"/>
 </tile>
 <tile id="109">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ae_v.png"/>
 </tile>
 <tile id="110">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_af_h.png"/>
 </tile>
 <tile id="111">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_af_v.png"/>
 </tile>
 <tile id="112">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_b8_h.png"/>
 </tile>
 <tile id="113">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_b8_v.png"/>
 </tile>
 <tile id="114">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_b9_h.png"/>
 </tile>
 <tile id="115">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_b9_v.png"/>
 </tile>
 <tile id="116">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ba_h.png"/>
 </tile>
 <tile id="117">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ba_v.png"/>
 </tile>
 <tile id="118">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_bb_h.png"/>
 </tile>
 <tile id="119">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_bb_v.png"/>
 </tile>
 <tile id="120">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_bc_h.png"/>
 </tile>
 <tile id="121">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_bc_v.png"/>
 </tile>
 <tile id="122">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_bd_h.png"/>
 </tile>
 <tile id="123">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_bd_v.png"/>
 </tile>
 <tile id="124">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_be_h.png"/>
 </tile>
 <tile id="125">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_be_v.png"/>
 </tile>
 <tile id="126">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_bf_h.png"/>
 </tile>
 <tile id="127">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_bf_v.png"/>
 </tile>
 <tile id="128">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_c8_h.png"/>
 </tile>
 <tile id="129">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_c8_v.png"/>
 </tile>
 <tile id="130">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_c9_h.png"/>
 </tile>
 <tile id="131">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_c9_v.png"/>
 </tile>
 <tile id="132">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ca_h.png"/>
 </tile>
 <tile id="133">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ca_v.png"/>
 </tile>
 <tile id="134">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_cb_h.png"/>
 </tile>
 <tile id="135">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_cb_v.png"/>
 </tile>
 <tile id="136">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_cc_h.png"/>
 </tile>
 <tile id="137">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_cc_v.png"/>
 </tile>
 <tile id="138">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_cd_h.png"/>
 </tile>
 <tile id="139">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_cd_v.png"/>
 </tile>
 <tile id="140">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ce_h.png"/>
 </tile>
 <tile id="141">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ce_v.png"/>
 </tile>
 <tile id="142">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_cf_h.png"/>
 </tile>
 <tile id="143">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_cf_v.png"/>
 </tile>
 <tile id="144">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_d8_h.png"/>
 </tile>
 <tile id="145">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_d8_v.png"/>
 </tile>
 <tile id="146">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_d9_h.png"/>
 </tile>
 <tile id="147">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_d9_v.png"/>
 </tile>
 <tile id="148">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_da_h.png"/>
 </tile>
 <tile id="149">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_da_v.png"/>
 </tile>
 <tile id="150">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_db_h.png"/>
 </tile>
 <tile id="151">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_db_v.png"/>
 </tile>
 <tile id="152">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_dc_h.png"/>
 </tile>
 <tile id="153">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_dc_v.png"/>
 </tile>
 <tile id="154">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_dd_h.png"/>
 </tile>
 <tile id="155">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_dd_v.png"/>
 </tile>
 <tile id="156">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_de_h.png"/>
 </tile>
 <tile id="157">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_de_v.png"/>
 </tile>
 <tile id="158">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_df_h.png"/>
 </tile>
 <tile id="159">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_df_v.png"/>
 </tile>
 <tile id="160">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_e8_h.png"/>
 </tile>
 <tile id="161">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_e8_v.png"/>
 </tile>
 <tile id="162">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_e9_h.png"/>
 </tile>
 <tile id="163">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_e9_v.png"/>
 </tile>
 <tile id="164">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ea_h.png"/>
 </tile>
 <tile id="165">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ea_v.png"/>
 </tile>
 <tile id="166">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_eb_h.png"/>
 </tile>
 <tile id="167">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_eb_v.png"/>
 </tile>
 <tile id="168">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ec_h.png"/>
 </tile>
 <tile id="169">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ec_v.png"/>
 </tile>
 <tile id="170">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ed_h.png"/>
 </tile>
 <tile id="171">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ed_v.png"/>
 </tile>
 <tile id="172">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ee_h.png"/>
 </tile>
 <tile id="173">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ee_v.png"/>
 </tile>
 <tile id="174">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ef_h.png"/>
 </tile>
 <tile id="175">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ef_v.png"/>
 </tile>
 <tile id="176">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_f8_h.png"/>
 </tile>
 <tile id="177">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_f8_v.png"/>
 </tile>
 <tile id="178">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_f9_h.png"/>
 </tile>
 <tile id="179">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_f9_v.png"/>
 </tile>
 <tile id="180">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_fa_h.png"/>
 </tile>
 <tile id="181">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_fa_v.png"/>
 </tile>
 <tile id="182">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_fb_h.png"/>
 </tile>
 <tile id="183">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_fb_v.png"/>
 </tile>
 <tile id="184">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_fc_h.png"/>
 </tile>
 <tile id="185">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_fc_v.png"/>
 </tile>
 <tile id="186">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_fd_h.png"/>
 </tile>
 <tile id="187">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_fd_v.png"/>
 </tile>
 <tile id="188">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_fe_h.png"/>
 </tile>
 <tile id="189">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_fe_v.png"/>
 </tile>
 <tile id="190">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ff_h.png"/>
 </tile>
 <tile id="191">
  <properties>
   <property name="opaque" type="bool" value="false"/>
   <property name="solid" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="bg_ff_v.png"/>
 </tile>
 <wangsets>
  <wangset name="Block/Open" tile="-1">
   <wangedgecolor name="Block" color="#000080" tile="-1" probability="1"/>
   <wangedgecolor name="Open" color="#ffff80" tile="-1" probability="1"/>
   <wangcornercolor name="Block" color="#000080" tile="-1" probability="1"/>
   <wangcornercolor name="Open" color="#ffff80" tile="-1" probability="1"/>
   <wangtile tileid="9" wangid="0x22222222"/>
   <wangtile tileid="10" wangid="0x21222222"/>
   <wangtile tileid="11" wangid="0x22212222"/>
   <wangtile tileid="12" wangid="0x21212222"/>
   <wangtile tileid="13" wangid="0x21112222"/>
   <wangtile tileid="14" wangid="0x22222122"/>
   <wangtile tileid="15" wangid="0x21222122"/>
   <wangtile tileid="16" wangid="0x22212122"/>
   <wangtile tileid="17" wangid="0x22211122"/>
   <wangtile tileid="18" wangid="0x21212122"/>
   <wangtile tileid="19" wangid="0x21112122"/>
   <wangtile tileid="20" wangid="0x21211122"/>
   <wangtile tileid="21" wangid="0x21111122"/>
   <wangtile tileid="22" wangid="0x22222221"/>
   <wangtile tileid="23" wangid="0x21222221"/>
   <wangtile tileid="24" wangid="0x11222221"/>
   <wangtile tileid="25" wangid="0x22212221"/>
   <wangtile tileid="26" wangid="0x21212221"/>
   <wangtile tileid="27" wangid="0x11212221"/>
   <wangtile tileid="28" wangid="0x21112221"/>
   <wangtile tileid="29" wangid="0x11112221"/>
   <wangtile tileid="30" wangid="0x22222121"/>
   <wangtile tileid="31" wangid="0x22222111"/>
   <wangtile tileid="32" wangid="0x21222121"/>
   <wangtile tileid="33" wangid="0x11222121"/>
   <wangtile tileid="34" wangid="0x21222111"/>
   <wangtile tileid="35" wangid="0x11222111"/>
   <wangtile tileid="36" wangid="0x22212121"/>
   <wangtile tileid="37" wangid="0x22211121"/>
   <wangtile tileid="38" wangid="0x22212111"/>
   <wangtile tileid="39" wangid="0x22211111"/>
   <wangtile tileid="40" wangid="0x21212121"/>
   <wangtile tileid="41" wangid="0x11212121"/>
   <wangtile tileid="42" wangid="0x21112121"/>
   <wangtile tileid="43" wangid="0x11112121"/>
   <wangtile tileid="44" wangid="0x21211121"/>
   <wangtile tileid="45" wangid="0x11211121"/>
   <wangtile tileid="46" wangid="0x21111121"/>
   <wangtile tileid="47" wangid="0x11111121"/>
   <wangtile tileid="48" wangid="0x21212111"/>
   <wangtile tileid="49" wangid="0x11212111"/>
   <wangtile tileid="50" wangid="0x21112111"/>
   <wangtile tileid="51" wangid="0x11112111"/>
   <wangtile tileid="52" wangid="0x21211111"/>
   <wangtile tileid="53" wangid="0x11211111"/>
   <wangtile tileid="54" wangid="0x21111111"/>
   <wangtile tileid="55" wangid="0x11111111"/>
  </wangset>
  <wangset name="Block/Open 1row" tile="-1">
   <wangedgecolor name="" color="#000080" tile="-1" probability="1"/>
   <wangedgecolor name="" color="#ffff80" tile="-1" probability="1"/>
   <wangtile tileid="9" wangid="0x2020202"/>
   <wangtile tileid="10" wangid="0x1020202"/>
   <wangtile tileid="11" wangid="0x2010202"/>
   <wangtile tileid="12" wangid="0x1010202"/>
   <wangtile tileid="14" wangid="0x2020102"/>
   <wangtile tileid="15" wangid="0x1020102"/>
   <wangtile tileid="16" wangid="0x2010102"/>
   <wangtile tileid="18" wangid="0x1010102"/>
   <wangtile tileid="22" wangid="0x2020201"/>
   <wangtile tileid="23" wangid="0x1020201"/>
   <wangtile tileid="25" wangid="0x2010201"/>
   <wangtile tileid="26" wangid="0x1010201"/>
   <wangtile tileid="30" wangid="0x2020101"/>
   <wangtile tileid="32" wangid="0x1020101"/>
   <wangtile tileid="36" wangid="0x2010101"/>
   <wangtile tileid="40" wangid="0x1010101"/>
  </wangset>
 </wangsets>
</tileset>
